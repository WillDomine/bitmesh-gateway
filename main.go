package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/WillDomine/bitmesh-gateway/internal/handler"
	"github.com/WillDomine/bitmesh-gateway/internal/router"
	"github.com/WillDomine/bitmesh-gateway/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := ":42000"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Initalizes the router in bitset_router.go and applies to grpcHandler 
	coreRouter := router.NewRouter()
	grpcHandler := &handler.GRPCHandler{
		Router: coreRouter,
	}

	//The GRPC server setup
	grpcServer := grpc.NewServer()
	pb.RegisterGatewayRouterServer(grpcServer, grpcHandler)

	//Debugging 
	reflection.Register(grpcServer)

	//Run server in Goroutine
	go func() {
		log.Printf("Bitmesh Gateway starting on %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
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