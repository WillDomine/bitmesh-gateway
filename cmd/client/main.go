package main

import (
	"context"
	"log"
	"time"

	"github.com/WillDomine/bitmesh-gateway/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//Connection to the gateway
	conn, err := grpc.NewClient("localhost:42000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGatewayRouterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//Standard User - Flag 0
	log.Println("Standard Request")
	req1, err := c.RouteMessage(ctx, &pb.RouteRequest{
		FeatureFlags: 0, 
		Payload: "Standard Request Test",
	})
	if err != nil {
		log.Printf("Standard request failed: %v", err)
	} else {
		log.Printf("Response: Service=%s, Allowed=%v", req1.TargetService, req1.Allowed)
	}

	//Premium User - Flag 8
	log.Println("Premium Request")
	req2, err := c.RouteMessage(ctx, &pb.RouteRequest{
		FeatureFlags: 8,
		Payload: "Premium Request Test",
	})
	if err != nil {
		log.Printf("Premium request failed: %v", err)
	} else {
		log.Printf("Response: Service=%s, Allowed=%v", req2.TargetService, req2.Allowed)
	}
}