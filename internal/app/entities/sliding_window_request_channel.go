package entities

import (
	"github.com/vaibhavahuja/rate-limiter/constants"
	"time"
)

type SlidingWindowRequestChannel struct {
	ServiceId       int
	CurrentTimeKey  string
	PreviousTimeKey string
	Duration        time.Duration
	CurrentTime     time.Time
	UnitType        constants.RateUnitType
	MaxLimit        int
	ResponseChannel chan SlidingWindowResponseChannel
}
