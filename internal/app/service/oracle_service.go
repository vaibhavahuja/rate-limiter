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
	log.Infof("my rate time unit is %v", rateTimeUnit)
	currentTime := time.Now()
	//creates key value
	currentTimeKeyValue, _ := json.Marshal(createKeyValue(int(serviceId), request, currentTime, 0, rateTimeUnit))
	previousTimeKeyValue, _ := json.Marshal(createKeyValue(int(serviceId), request, currentTime, rateTimeUnit.GetDuration(), rateTimeUnit))
	log.Infof("here are my keys : currentTimeKey : %s, previousTimeKeyValue : %s", currentTimeKeyValue, previousTimeKeyValue)
	//fetches value of prevKey and currKey counters concurrently!!

	//todo fetch both in one call to cache
	currentCounter, _ := app.requestCounterCache.FetchCounterValueForKey(string(currentTimeKeyValue))
	prevCounter, _ := app.requestCounterCache.FetchCounterValueForKey(string(previousTimeKeyValue))

	exists := false
	log.Infof("currentCounter : %v, prevCounter %v", currentCounter, prevCounter)
	if currentCounter != 0 {
		exists = true
	}
	//your logic - sliding window algorithm
	//todo write sliding window algo here

	//send update request to channel -> do this asynchronously
	err = app.requestCounterCache.IncrementRequestCounter(string(currentTimeKeyValue), exists, 2*rateTimeUnit.GetDuration())
	if err != nil {
		log.Errorf("error while incrementing request counter")
	}
	return true, nil
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