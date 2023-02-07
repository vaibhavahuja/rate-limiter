package constants

type RateUnitType int

const (
	Second RateUnitType = iota
	Minute
	Hour
	Day
	Month
)

func ValidateRateUnitType(value int) bool {
	if v := RateUnitType(value); v > Month || v < Second {
		return false
	}
	return true
}
