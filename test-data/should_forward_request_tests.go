package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	pb "github.com/vaibhavahuja/rate-limiter/proto"
	"testing"
)

func ShouldForwardRequest(t *testing.T, client pb.RateLimiterClient) {
	shouldForwardRequests(t, client)
}

func shouldForwardRequests(t *testing.T, client pb.RateLimiterClient) {
	ctx := context.TODO()
	input := &pb.ShouldForwardsRequestRequest{
		ServiceId: 1,
		Request:   "{test_request}",
	}
	requestsPassed := 0
	for i := 0; i < 15; i++ {
		resp, _ := client.ShouldForwardRequest(ctx, input)
		if resp.GetShouldForward() {
			requestsPassed++
		}
	}
	//service rate limit was defined as 10 requests per minute
	assert.Equal(t, 10, requestsPassed)
}
