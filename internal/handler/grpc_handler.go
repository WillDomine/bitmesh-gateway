package handler

import (
	"context"
	"fmt"

	"github.com/WillDomine/bitmesh-gateway/api/pb"
	"github.com/WillDomine/bitmesh-gateway/internal/proxy"
	"github.com/WillDomine/bitmesh-gateway/internal/router"
)

type GRPCHandler struct {
	pb.UnimplementedGatewayRouterServer
	Router    *router.Router
	Forwarder *proxy.Forwarder
}

func (h *GRPCHandler) RouteMessage(ctx context.Context, req *pb.RouteRequest) (*pb.RouteResponse, error) {
	// The routing happens here using the bitwise logic
	target, allowed := h.Router.Decide(req.FeatureFlags)
	if !allowed {
		return nil, fmt.Errorf("permission denied")
	}

	//Map the service address to var (Future I will add a config.yaml that handles this)
	var targetAddr string
	switch target {
	case "premium-backend":
		targetAddr = "localhost:50000"
	case "standard-backend":
		targetAddr = "localhost:50001"
	default:
		return nil, fmt.Errorf("unknown service: %s", target)
	}

	//Send the traffic to the address thorugh forward func in proxy/forwarder.go
	return h.Forwarder.Forward(ctx, targetAddr, req)
}
