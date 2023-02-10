package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/vaibhavahuja/rate-limiter/constants"
	"testing"
	"time"
)

func TestGetSlidingWindowRequestCount(t *testing.T) {
	currentTime := time.Date(2022, 12, 1, 12, 30, 50, 0, time.UTC)
	currentWindowCounter := 10
	previousWindowCounter := 50
	ans := GetSlidingWindowRequestCount(currentTime,constants.RateUnitType(2), currentWindowCounter, previousWindowCounter)
	assert.Equal(t, ans, 34)
}
