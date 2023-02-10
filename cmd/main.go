package main

import (
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/config"
	"github.com/vaibhavahuja/rate-limiter/constants"
	"github.com/vaibhavahuja/rate-limiter/internal/app/entities"
	"github.com/vaibhavahuja/rate-limiter/internal/app/gateway/cache"
	"github.com/vaibhavahuja/rate-limiter/internal/app/gateway/handler"
	"github.com/vaibhavahuja/rate-limiter/internal/app/gateway/repository"
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
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	log.SetReportCaller(true)
}

func main() {
	log.Info("Starting Rate Limiter service")
	s := startGRPCServer()
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
	dynamoClient := infrastructure.GetDynamoDBClient(conf.RateLimiterDynamo.Region, conf.RateLimiterDynamo.Endpoint)
	rulesRepository := repository.NewRuleRepository(dynamoClient, conf.RateLimiterDynamo.Table)
	requestCounterCache := cache.NewRequestCounterCache(redisClient)
	rateLimitingRequestChanArray := initialiseRateLimiterChannels(constants.RateLimitingRequestWorkersSize, constants.RateLimitingChannelQueueSize)
	exitChan := make(chan bool)
	application := service.NewApplication(conf, rulesRepository, requestCounterCache)
	log.Infof("after initialising length of channel : %d", len(rateLimitingRequestChanArray))
	//starting worker goRoutines
	for i := 0; i < constants.RateLimitingRequestWorkersSize; i++ {
		go application.SlidingWindowAlgorithmWorker(i, rateLimitingRequestChanArray[i], exitChan)
	}
	svc := handler.NewRateLimiterGrpcServer(application, rateLimitingRequestChanArray)
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

func initialiseRateLimiterChannels(workerSize, queueSize int) (resp []chan entities.SlidingWindowRequestChannel) {
	//initialising requestChannelArray
	for i := 0; i < workerSize; i++ {
		resp = append(resp, make(chan entities.SlidingWindowRequestChannel, queueSize))
	}
	return
}

func gracefulStop(s *grpc.Server) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal := <-c
	log.Infof("Stopping rate limiter service. signal : %s", signal)
	s.GracefulStop()
}
