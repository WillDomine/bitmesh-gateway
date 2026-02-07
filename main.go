package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/WillDomine/bitmesh-gateway/api/pb"
	"github.com/WillDomine/bitmesh-gateway/internal/config"
	"github.com/WillDomine/bitmesh-gateway/internal/handler"
	"github.com/WillDomine/bitmesh-gateway/internal/proxy"
	"github.com/WillDomine/bitmesh-gateway/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	//Allows path to be overridden via CLI flags (mismatch between location for manual run vs docker)
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()
	//Loads the configuration file (server port and services to call)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config, %v", err)
	}

	//Links the server of type tcp to port
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//Initalizes the router in bitset_router.go and applies to grpcHandler
	coreRouter := router.NewRouter()
	netForwarder := proxy.NewForwarder()

	//handles all core components
	grpcHandler := &handler.GRPCHandler{
		Router:     coreRouter,
		Forwarder:  netForwarder,
		ServiceMap: cfg.Services,
	}

	//The GRPC server setup
	grpcServer := grpc.NewServer()
	pb.RegisterGatewayRouterServer(grpcServer, grpcHandler)

	//Debugging
	reflection.Register(grpcServer)

	//Run server in Goroutine
	go func() {
		log.Printf("Bitmesh Gateway starting on %s", cfg.Server.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	//Graceful shutdown waiting for interrupt signal
	stopChain := make(chan os.Signal, 1)
	signal.Notify(stopChain, os.Interrupt, syscall.SIGTERM)

	<-stopChain
	log.Printf("Shutting down server")

	grpcServer.GracefulStop()
	log.Printf("Server stopped")
}
