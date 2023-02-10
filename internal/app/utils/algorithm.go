package utils

import (
	"github.com/vaibhavahuja/rate-limiter/constants"
	"time"
)

// GetSlidingWindowRequestCount Returns the count of requests for the specified timeUnit
func GetSlidingWindowRequestCount(currentTime time.Time, rateUnitType constants.RateUnitType, currentWindowCounter, previousWindowCounter int) (count int) {
	percentageTimePassed := rateUnitType.GetTimePassedPercentage(currentTime.UTC())
	count = int((1-percentageTimePassed)*float64(previousWindowCounter) + float64(currentWindowCounter))
	return
}
