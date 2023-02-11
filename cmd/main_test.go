package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/config"
	pb "github.com/vaibhavahuja/rate-limiter/proto"
	tests "github.com/vaibhavahuja/rate-limiter/test-data"
	"github.com/vaibhavahuja/rate-limiter/test-data/dynamo"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestFunctional(t *testing.T) {
	log.Infof("Running functional tests")
	sltConfig := config.GetConfig()
	err := InitTestData(sltConfig)
	if err != nil {
		log.Errorf("error whiles initialising test data %s", err)
		return
	}
	go main()
	time.Sleep(2 * time.Second)
	grpcClient := initGrpcClient(sltConfig)
	tests.RegisterServiceTests(t, grpcClient)
	tests.ShouldForwardRequest(t, grpcClient)
}

func InitTestData(sltConfig *config.Config) error {
	log.Infof("initialising for functional tests")
	//creating dynamoDB service table
	err := dynamo.CreateServiceTable(sltConfig)
	//either error can be that table already exists or it can be something else?
	if err != nil {
		log.Errorf("error while creating service table / table already exists")
		return nil
	}
	return nil
	// initialising redis? I mean that can be done via the up file as well?
}

func initGrpcClient(conf *config.Config) pb.RateLimiterClient {
	target := "localhost:" + conf.Service.Port
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Errorf("unable to dial target %s error : %s", target, err.Error())
	}
	rateLimiterClient := pb.NewRateLimiterClient(conn)
	return rateLimiterClient
}
