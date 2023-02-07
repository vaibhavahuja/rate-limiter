package service

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/vaibhavahuja/rate-limiter/constants"
	"github.com/vaibhavahuja/rate-limiter/internal/app/entities"
	"github.com/vaibhavahuja/rate-limiter/proto"
)

// RegisterService saves the rule and service id to datastore
func (app *Application) RegisterService(ctx context.Context, serviceId int, rule *proto.Rule) error {
	if serviceId == constants.DefaultServiceId {
		return errors.New("invalid service id")
	}
	if !constants.ValidateRateUnitType(int(rule.GetRate().GetUnit())) {
		return errors.New("invalid rate unit type")
	}

	log.Infof("Registering service with id %d", serviceId)

	ruleInput := entities.Rule{
		Field: rule.Field,
		Rate: entities.Rate{
			RequestsPerUnit: rule.Rate.RequestsPerUnit,
			UnitType:        rule.Rate.Unit,
		},
		RequestRejectionMessage: rule.ErrorMessage,
	}
	err := app.rulesRepository.AddService(ctx, serviceId, ruleInput)
	if err != nil {
		return err
	}
	log.Infof("successfully registered service: %d", serviceId)

	return nil
}
