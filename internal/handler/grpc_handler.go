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
	Router     *router.Router
	Forwarder  *proxy.Forwarder
	ServiceMap map[string]string
}

func (h *GRPCHandler) RouteMessage(ctx context.Context, req *pb.RouteRequest) (*pb.RouteResponse, error) {
	// The routing happens here using the bitwise logic
	target, allowed := h.Router.Decide(req.FeatureFlags)
	if !allowed {
		return nil, fmt.Errorf("permission denied")
	}

	//Search the map registry for the service
	targetAddr, exitsts := h.ServiceMap[target]
	if !exitsts {
		return nil, fmt.Errorf("Service not found: %s", target)
	}
	//Send the traffic to the address thorugh forward func in proxy/forwarder.go
	return h.Forwarder.Forward(ctx, targetAddr, req)
}
