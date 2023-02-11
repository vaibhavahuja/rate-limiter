package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	pb "github.com/vaibhavahuja/rate-limiter/proto"
	"testing"
)

func RegisterServiceTests(t *testing.T, client pb.RateLimiterClient) {
	registerServiceTestSuccess(t, client)
	registerServiceTestFailureInvalidRateUnit(t, client)
	registerServiceTestFailureInvalidServiceId(t, client)
}

func registerServiceTestSuccess(t *testing.T, client pb.RateLimiterClient) {
	ctx := context.TODO()
	input := &pb.RegisterServiceRequest{
		ServiceId: 1,
		Rule: &pb.Rule{
			Field: "",
			Rate: &pb.Rate{
				RequestsPerUnit: 10,
				Unit:            1,
			},
			ErrorMessage: "retry-with-exponential-backoff",
		},
	}
	resp, err := client.RegisterService(ctx, input)
	if err != nil {
		t.Errorf("error while registering service %s", err.Error())
	}

	assert.Equal(t, int32(1), resp.GetStatus())
}

func registerServiceTestFailureInvalidRateUnit(t *testing.T, client pb.RateLimiterClient) {
	ctx := context.TODO()
	input := &pb.RegisterServiceRequest{
		ServiceId: 1,
		Rule: &pb.Rule{
			Field: "",
			Rate: &pb.Rate{
				RequestsPerUnit: 10,
				Unit:            8,
			},
			ErrorMessage: "retry-with-exponential-backoff",
		},
	}
	resp, err := client.RegisterService(ctx, input)
	expectedError := "rpc error: code = Unknown desc = invalid rate unit type"
	assert.Equal(t, expectedError, err.Error())
	assert.Equal(t, int32(0), resp.GetStatus())
}

func registerServiceTestFailureInvalidServiceId(t *testing.T, client pb.RateLimiterClient) {
	ctx := context.TODO()
	input := &pb.RegisterServiceRequest{
		ServiceId: 0,
		Rule: &pb.Rule{
			Field: "",
			Rate: &pb.Rate{
				RequestsPerUnit: 10,
				Unit:            8,
			},
			ErrorMessage: "retry-with-exponential-backoff",
		},
	}
	resp, err := client.RegisterService(ctx, input)
	expectedError := "rpc error: code = Unknown desc = invalid service id"
	assert.Equal(t, expectedError, err.Error())
	assert.Equal(t, int32(0), resp.GetStatus())
}
