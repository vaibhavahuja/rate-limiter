package service

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/constants"
	"github.com/vaibhavahuja/rate-limiter/internal/app/entities"
	"github.com/vaibhavahuja/rate-limiter/internal/app/utils"
	"time"
)

// ShouldForwardRequest returns true if the rate limit of serviceId has not reached
func (app *Application) ShouldForwardRequest(ctx context.Context, serviceId int32, request string, slidingWindowRequestChannel []chan entities.SlidingWindowRequestChannel) (bool, error) {
	log.Infof("Checking if rate limit has reached or not for : %d", serviceId)
	//fetching the rule for given serviceId
	serviceRule, err := app.rulesRepository.GetRuleByServiceId(ctx, int(serviceId))
	if err != nil {
		log.Errorf("error while fetching rules for service id %d", serviceId)
		return true, err
	}
	rateTimeUnit := constants.RateUnitType(serviceRule.Rate.UnitType)
	currentTime := time.Now()
	currentTimeKeyValue, _ := json.Marshal(createKeyValue(int(serviceId), request, currentTime, 0, rateTimeUnit))
	previousTimeKeyValue, _ := json.Marshal(createKeyValue(int(serviceId), request, currentTime, rateTimeUnit.GetDuration(), rateTimeUnit))
	responseChannel := make(chan entities.SlidingWindowResponseChannel)
	//forwarding requests of all service_id x to the same channel
	slidingWindowRequestChannel[serviceId%constants.RateLimitingRequestWorkersSize] <- entities.SlidingWindowRequestChannel{
		CurrentTimeKey:  string(currentTimeKeyValue),
		PreviousTimeKey: string(previousTimeKeyValue),
		Duration:        2 * rateTimeUnit.GetDuration(),
		CurrentTime:     currentTime,
		MaxLimit:        int(serviceRule.Rate.RequestsPerUnit),
		UnitType:        rateTimeUnit,
		ServiceId:       int(serviceId),
		ResponseChannel: responseChannel,
	}
	response := <-responseChannel
	return response.ShouldForward, err
}

// SlidingWindowAlgorithmWorker reads current and previous key and increments currentKey
// would need multiple go routines ready each listening on separate chanel
func (app *Application) SlidingWindowAlgorithmWorker(id int, slidingWindowRequestChannel <-chan entities.SlidingWindowRequestChannel, exitChan <-chan bool) {
	log.Infof("starting new sliding window worker with id %d", id)
	for {
		select {
		case info := <-slidingWindowRequestChannel:
			log.Infof("received request on request channel(%d) %v", id, info)
			counterVals, _ := app.requestCounterCache.FetchCounterValueForKeys(info.CurrentTimeKey, info.PreviousTimeKey)
			currentCounter, prevCounter := counterVals[0], counterVals[1]
			exists := false
			if currentCounter != 0 {
				exists = true
			}
			totalRequestsCountInSlidingWindow := utils.GetSlidingWindowRequestCount(info.CurrentTime, info.UnitType, currentCounter, prevCounter)
			if totalRequestsCountInSlidingWindow >= info.MaxLimit {
				log.Infof("rejecting the request for service %d", info.ServiceId)
				info.ResponseChannel <- entities.SlidingWindowResponseChannel{
					ServiceId:     info.ServiceId,
					ShouldForward: false,
				}
				continue
			}
			//sending response to response channel before incrementing request counter in redis
			info.ResponseChannel <- entities.SlidingWindowResponseChannel{
				ServiceId:     info.ServiceId,
				ShouldForward: true,
			}
			err := app.requestCounterCache.IncrementRequestCounter(info.CurrentTimeKey, exists, info.Duration)
			if err != nil {
				log.Errorf("error while incrementing request counter %v", err.Error())
			}
		case <-exitChan:
			return
		}
	}
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
