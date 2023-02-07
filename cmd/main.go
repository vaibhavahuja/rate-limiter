package main

import (
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/config"
	"github.com/vaibhavahuja/rate-limiter/internal/app/gateway"
	"github.com/vaibhavahuja/rate-limiter/internal/app/service"
	"github.com/vaibhavahuja/rate-limiter/internal/pkg/infrastructure"
	pb "github.com/vaibhavahuja/rate-limiter/proto"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	conf *config.Config
)

func init() {
	conf = config.GetConfig()
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	log.SetReportCaller(true)
}

func main() {
	log.Info("starting rate limiter")
	s := startGRPCServer()

	log.Info("started rate limiter")
	gracefulStop(s)
}

func startGRPCServer() *grpc.Server {
	lis, err := net.Listen("tcp", ":"+conf.Service.Port)
	if err != nil {
		log.Errorf("Error while starting GRPC Server %s", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				//metrics.Prometheus.Grpc.ServerMetrics.UnaryServerInterceptor(),
				grpcrecovery.UnaryServerInterceptor(),
			),
		),
	)
	redisClient := infrastructure.GetRedisClient(conf.Redis.Url, conf.Redis.Password)
	application := service.NewApplication(conf, redisClient)
	svc := gateway.NewRateLimiterGrpcServer(application)

	pb.RegisterRateLimiterServer(s, &svc)
	go func() {
		log.Info("Starting gRPC server")
		if err := s.Serve(lis); err != nil {
			log.Infof("failed to serve %s", err)
		}
	}()
	log.Infof("Successfully started gRPC server and running on %v", s.GetServiceInfo())
	return s
}

func gracefulStop(s *grpc.Server) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal := <-c
	log.Infof("Stopping rate limiter service. signal : %s", signal)
	s.GracefulStop()
}
