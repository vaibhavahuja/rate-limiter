package service

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/constants"
	"github.com/vaibhavahuja/rate-limiter/internal/app/entities"
	"time"
)

// ShouldForwardRequest returns true if the rate limit of serviceId has not reached
func (app *Application) ShouldForwardRequest(ctx context.Context, serviceId int32, request string) (bool, error) {
	log.Infof("Checking if rate limit has reached or not for : %d", serviceId)
	//fetches the rule for given serviceId
	serviceRule, err := app.rulesRepository.GetRuleByServiceId(ctx, int(serviceId))
	if err != nil {
		log.Errorf("error while fetching rules for service id %d", serviceId)
		return true, err
	}
	rateTimeUnit := constants.RateUnitType(serviceRule.Rate.UnitType)
	currentTime := time.Now()
	//creates key value
	currentTimeKeyValue, _ := json.Marshal(createKeyValue(int(serviceId), request, currentTime, 0, rateTimeUnit))
	previousTimeKeyValue, _ := json.Marshal(createKeyValue(int(serviceId), request, currentTime, rateTimeUnit.GetDuration(), rateTimeUnit))
	log.Debugf("here are my keys : currentTimeKey : %s, previousTimeKeyValue : %s", currentTimeKeyValue, previousTimeKeyValue)

	//todo fetch both in one call to cache
	currentCounter, _ := app.requestCounterCache.FetchCounterValueForKey(string(currentTimeKeyValue))
	prevCounter, _ := app.requestCounterCache.FetchCounterValueForKey(string(previousTimeKeyValue))

	currentCounterExists := false
	log.Debugf("currentCounter : %v, prevCounter %v", currentCounter, prevCounter)

	if currentCounter != 0 {
		currentCounterExists = true
	}
	log.Infof("max limit is : %v and you have reached %v requests", serviceRule.Rate.RequestsPerUnit, getSlidingWindowRequestCount(currentTime, rateTimeUnit, currentCounter, prevCounter))
	totalRequestsCountInSlidingWindow := getSlidingWindowRequestCount(currentTime, rateTimeUnit, currentCounter, prevCounter)
	if totalRequestsCountInSlidingWindow >= int(serviceRule.Rate.RequestsPerUnit) {
		//We do not allow the request to go through
		return false, nil
	}
	//Spinning new goRoutine to increment request counter
	go func(key string, exists bool, ttl time.Duration) {
		//Register request if the request goes through
		err = app.requestCounterCache.IncrementRequestCounter(string(currentTimeKeyValue), currentCounterExists, 2*rateTimeUnit.GetDuration())
		if err != nil {
			log.Errorf("error while incrementing request counter")
		}
	}(string(currentTimeKeyValue), currentCounterExists, 2*rateTimeUnit.GetDuration())

	return true, nil
}

// slidingWindowRequestCounter Returns the count of requests for the specified timeUnit
func getSlidingWindowRequestCount(currentTime time.Time, rateUnitType constants.RateUnitType, currentWindowCounter, previousWindowCounter int) (count int) {
	percentageTimePassed := rateUnitType.GetTimePassedPercentage(currentTime.UTC())
	//log.Debugf("calculated percentage time passed : %v", percentageTimePassed)
	count = int((1-percentageTimePassed)*float64(previousWindowCounter) + float64(currentWindowCounter))
	return
}

// createKeyValue Fetches keyValue for storing/fetching counter from cache
func createKeyValue(serviceId int, field string, currentTime time.Time, subtractDuration time.Duration, timeUnit constants.RateUnitType) entities.RequestCounterCacheKey {
	subtractedTime := currentTime.Add(-1 * subtractDuration)
	log.Debugf("received request to create key value currentTime : %v,  subtracted duration: %v,", currentTime, subtractDuration)
	log.Debugf("converted time to rate unit %v", timeUnit.ConvertTimeToRateUnit(subtractedTime))
	return entities.RequestCounterCacheKey{
		ServiceId: serviceId,
		Field:     field,
		TimeValue: timeUnit.ConvertTimeToRateUnit(subtractedTime),
	}
}
