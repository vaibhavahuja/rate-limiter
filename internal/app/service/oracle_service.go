package service

import (
	log "github.com/sirupsen/logrus"
)

//ShouldForwardRequest returns true if the rate limit of serviceId has not reached
func (app *Application) ShouldForwardRequest(serviceId int32, request string) (bool, error) {
	log.Infof("Checking if rate limit has reached or not for : %d", serviceId)
	//todo implement the method
	//fetch rules + fetch data from cache
	//asynchronously update cache
	//return response
	return false, nil
}
