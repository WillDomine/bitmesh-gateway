package proxy

import (
	"context"
	"fmt"
	"log"

	"github.com/WillDomine/bitmesh-gateway/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//Normally would create connection pool so we don't dial on every request, but for now I will dial on each request
type Forwarder struct {}

func NewForwarder() *Forwarder {
	return &Forwarder{}
}

//Sends the request to the targeted service and returns the response
func (f *Forwarder) Forward(ctx context.Context, targetAddr string, req *pb.RouteRequest) (*pb.RouteResponse, error) {
	log.Printf("Proxy request: %s", targetAddr)

	//Connects to service and using insecure for simplicity since no TLS between microservices
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service: %w", err)
	}
	defer conn.Close()

	client := pb.NewGatewayRouterClient(conn)

	resp, err := client.RouteMessage(ctx, req) 
	if err != nil {
		return nil, fmt.Errorf("service failed: %w", err)
	}

	return resp, nil
}