package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/WillDomine/bitmesh-gateway/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//Testing Config
	const (
		RequestsToSend = 10000 //Requests
		Concurrency    = 50    //Users
		GatewayAddr    = "localhost:50000"
	)

	log.Printf("Starting Load Test: #%d requests with #%d users", RequestsToSend, Concurrency)

	//Set up clients
	clients := make([]pb.GatewayRouterClient, Concurrency)
	for i := 0; i < Concurrency; i++ {
		conn, err := grpc.NewClient(GatewayAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		clients[i] = pb.NewGatewayRouterClient(conn)
	}

	//Metrics
	var (
		successCount int64
		failCount    int64
		totalLatency int64
	)

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(Concurrency)

	requestsPerUser := RequestsToSend / Concurrency

	// Launch Users
	for i := 0; i < Concurrency; i++ {
		go func(workerID int) {
			defer wg.Done()
			client := clients[workerID]

			for j := 0; j < requestsPerUser; j++ {
				startReq := time.Now()

				// Alternate between Standard (0) and Premium (8)
				flags := uint32(0)
				if j%2 == 0 {
					flags = 8
				}

				_, err := client.RouteMessage(context.Background(), &pb.RouteRequest{
					FeatureFlags: flags,
					Payload:      "bench-load",
				})

				duration := time.Since(startReq).Microseconds()
				atomic.AddInt64(&totalLatency, duration)

				if err != nil {
					atomic.AddInt64(&failCount, 1)
				} else {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	rps := float64(successCount) / elapsed.Seconds()
	avgLatency := float64(totalLatency) / float64(successCount) / 1000.0 //mili seconds

	fmt.Println("\n--- BitMesh Performance Metrics ---")
	fmt.Printf("Total Requests: %d\n", RequestsToSend)
	fmt.Printf("Concurrency:    %d workers\n", Concurrency)
	fmt.Printf("Time Taken:     %v\n", elapsed)
	fmt.Printf("Success/Fail:   %d / %d\n", successCount, failCount)
	fmt.Println("-----------------------------------")
	fmt.Printf("Throughput:     %.2f req/sec\n", rps)
	fmt.Printf("Avg Latency:    %.4f ms\n", avgLatency)
	fmt.Println("-----------------------------------")
}
