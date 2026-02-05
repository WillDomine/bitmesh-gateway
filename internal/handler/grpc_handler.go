package handler

import (
	"context"
	"github.com/WillDomine/bitmesh-gateway/api/pb"
	"github.com/WillDomine/bitmesh-gateway/internal/router"
)

type GRPCHandler struct {
	pb.UnimplementedGatewayRouterServer
	Router *router.Router
}

func (h *GRPCHandler) RouteMessage(ctx context.Context, req *pb.RouteRequest) (*pb.RouteResponse, error) {
	// The routing happens here using the bitwise logic
	target, allowed := h.Router.Decide(req.FeatureFlags)
	
	return &pb.RouteResponse{
		TargetService: target,
		Allowed:       allowed,
	}, nil
}