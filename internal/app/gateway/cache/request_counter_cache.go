package cache

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type RequestCounterCache Cache

func NewRequestCounterCache(client *redis.Client) *RequestCounterCache {
	return &RequestCounterCache{
		redisClient: client,
	}
}

// IncrementRequestCounter increments request counter by one if key exists
// otherwise, it creates the key with defined TTL
func (rcc *RequestCounterCache) IncrementRequestCounter(key string, exists bool, ttl time.Duration) (err error) {
	if exists {
		err = rcc.redisClient.Incr(key).Err()
		if err != nil {
			log.Errorf("error while incrementing request counter for key : %s", key)
			return
		}
		return
	}
	//create key
	err = rcc.redisClient.Set(key, 1, ttl).Err()
	if err != nil {
		log.Errorf("error while setting key : %s", key)
		return
	}
	return
}

// FetchCounterValueForKey fetches the value for counters of given key
func (rcc *RequestCounterCache) FetchCounterValueForKey(key string) (int, error) {
	counterVal, err := rcc.redisClient.Get(key).Result()
	if err != nil {
		log.Errorf("Error while fetching counter value for key : %s", key)
		return 0, err
	}
	intVal, _ := strconv.Atoi(counterVal)
	return intVal, nil
}
