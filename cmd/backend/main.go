package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/WillDomine/bitmesh-gateway/api/pb"
	"google.golang.org/grpc"
)

type DummyTestingServer struct {
	pb.UnimplementedGatewayRouterServer
	Name string
}

// Routes the messages from the router service
func (d *DummyTestingServer) RouteMessage(ctx context.Context, req *pb.RouteRequest) (*pb.RouteResponse, error) {
	log.Printf("%s has received a request with payload %s", d.Name, req.Payload)
	return &pb.RouteResponse{
		TargetService: d.Name,
		Allowed:       true,
	}, nil
}

// Starts the dummy server
func startServer(name, port string, wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to list on %s: %v", port, err)
	}

	server := grpc.NewServer()
	pb.RegisterGatewayRouterServer(server, &DummyTestingServer{Name: name})

	log.Printf("Starting %s on %s", name, port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve %s: %v", name, err)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	//Premium Server
	go startServer("premium-backend", ":50000", &wg)

	//Standard Server
	go startServer("standard-backend", ":50001", &wg)

	wg.Wait()
}
