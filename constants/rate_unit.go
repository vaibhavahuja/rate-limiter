package constants

import (
	"time"
)

type RateUnitType int

const (
	Second RateUnitType = iota
	Minute
	Hour
	Day
)

func ValidateRateUnitType(value int) bool {
	if v := RateUnitType(value); v > Day || v < Second {
		return false
	}
	return true
}

// GetDurationByRateUnitTypeId gets the time duration by RateUnitTypeId
func (rt RateUnitType) GetDuration() time.Duration {
	switch rt {
	case Second:
		return time.Second
	case Minute:
		return time.Minute
	case Hour:
		return time.Hour
	case Day:
		return 24 * time.Hour
	default:
		return time.Minute
	}
}

// ConvertTimeToRateUnit converts the time upto specified rateUnit
func (rt RateUnitType) ConvertTimeToRateUnit(input time.Time) time.Time {
	switch rt {
	case Second:
		return time.Date(input.Year(), input.Month(), input.Day(), input.Hour(), input.Minute(), input.Second(), 0, time.UTC)
	case Minute:
		return time.Date(input.Year(), input.Month(), input.Day(), input.Hour(), input.Minute(), 0, 0, time.UTC)
	case Hour:
		return time.Date(input.Year(), input.Month(), input.Day(), input.Hour(), 0, 0, 0, time.UTC)
	case Day:
		return time.Date(input.Year(), input.Month(), input.Day(), 0, 0, 0, 0, time.UTC)
	default:
		return time.Date(input.Year(), input.Month(), input.Day(), input.Hour(), input.Minute(), input.Second(), input.Nanosecond(), time.UTC)
	}
}
