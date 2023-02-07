package gateway

import (
	"context"
	"github.com/vaibhavahuja/rate-limiter/internal/app/service"
	pb "github.com/vaibhavahuja/rate-limiter/proto"
)

type RateLimiterGrpcServer struct {
	service *service.Application
}

func NewRateLimiterGrpcServer(service *service.Application) RateLimiterGrpcServer {
	return RateLimiterGrpcServer{service: service}
}


func (rlServer *RateLimiterGrpcServer) mustEmbedUnimplementedRateLimiterServer() {
	//TODO implement me
	panic("implement me")
}


func (rlServer *RateLimiterGrpcServer) RegisterService(ctx context.Context, req *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	//todo to implement this
	return nil, nil
}

func (rlServer *RateLimiterGrpcServer) ShouldForwardRequest(ctx context.Context, in *pb.ShouldForwardsRequestRequest) (*pb.ShouldForwardRequestResponse, error) {
	//todo to implement this
	return nil, nil
}
