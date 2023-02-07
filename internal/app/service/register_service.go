package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/proto"
)

//RegisterService saves the rule and service id to datastore
func (app *Application) RegisterService(serviceId int32, rule *proto.Rule) error {
	log.Infof("Registering service with id %d", serviceId)
	//todo implement the method
	return nil
}
