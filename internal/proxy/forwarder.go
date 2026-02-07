package proxy

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/WillDomine/bitmesh-gateway/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//Maintains a connection pool of gRPC connections to reduce latency
type Forwarder struct {
	//Mutex to unlock and lock, so people don't add the same connections
	mutex sync.RWMutex
	//The map that stores the conn using the addr as a key
	connections map[string]*grpc.ClientConn
}

func NewForwarder() *Forwarder {
	return &Forwarder{
		//Allocates space for client map
		connections: make(map[string]*grpc.ClientConn),
	}
}

// Sends the request to the targeted service and returns the response
func (f *Forwarder) Forward(ctx context.Context, targetAddr string, req *pb.RouteRequest) (*pb.RouteResponse, error) {
	log.Printf("Proxy request: %s", targetAddr)

	//Check if the connection exists already
	f.mutex.RLock()
	conn, exists := f.connections[targetAddr]
	f.mutex.RUnlock()

	//Double Check Safety
	//Request the targeted service and add them to the pool
	if !exists {
		f.mutex.Lock()
		//Check again incase it was pooled it
		conn, exists = f.connections[targetAddr]
		if !exists {
			log.Printf("Pool Missing Service Requesting: %s", targetAddr)
			//Allocate for new client err and pre-allocated conn
			var err error
			conn, err = grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				//No Deadlocks here
				f.mutex.Unlock()
				return nil, fmt.Errorf("failed to connect to service: %w", err)
			}
			f.connections[targetAddr] = conn
		}
		f.mutex.Unlock()
	}

	//Create the client with the connection
	client := pb.NewGatewayRouterClient(conn)

	//Request the service with the client
	resp, err := client.RouteMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("service failed: %w", err)
	}
	//Returns the response and no errors
	return resp, nil
}
