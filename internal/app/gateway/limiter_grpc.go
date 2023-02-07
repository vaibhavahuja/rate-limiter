package gateway

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/internal/app/service"
	pb "github.com/vaibhavahuja/rate-limiter/proto"
)

type RateLimiterGrpcServer struct {
	service *service.Application
}

func NewRateLimiterGrpcServer(service *service.Application) RateLimiterGrpcServer {
	return RateLimiterGrpcServer{service: service}
}

func (rlServer *RateLimiterGrpcServer) RegisterService(ctx context.Context, req *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	log.Infof("received request for register service : %v", req)
	//todo add metrics for latency and request count
	err := rlServer.service.RegisterService(req.GetServiceId(), req.GetRule())
	if err != nil {
		log.Errorf("unable to register service. here's why : %s", err.Error())
		return nil, err
	}
	//todo should I just use google status codes?
	return &pb.RegisterServiceResponse{
		Status: 1,
		Error:  "",
	}, nil
}

func (rlServer *RateLimiterGrpcServer) ShouldForwardRequest(ctx context.Context, req *pb.ShouldForwardsRequestRequest) (*pb.ShouldForwardRequestResponse, error) {
	log.Infof("received request for register service : %v", req)
	//todo add metrics for latency and request count
	shouldForward, err := rlServer.service.ShouldForwardRequest(req.GetServiceId(), req.GetRequest())
	if err != nil {
		log.Errorf("error while seeing if rate limit has reached or not. here's why : %s", err.Error())
	}
	//todo should I just use google status codes?
	return &pb.ShouldForwardRequestResponse{
		ShouldForward: shouldForward,
	}, nil
}
