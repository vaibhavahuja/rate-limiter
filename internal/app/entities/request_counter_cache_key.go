package entities

import "time"

type RequestCounterCacheKey struct {
	ServiceId int
	Field     string
	TimeValue time.Time
}
