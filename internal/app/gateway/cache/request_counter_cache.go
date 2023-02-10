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

// IncrementRequestCounter Increments the request counter and sets the expiry if key is getting initialised
func (rcc *RequestCounterCache) IncrementRequestCounter(key string, exists bool, ttl time.Duration) (err error) {
	err = rcc.redisClient.Incr(key).Err()
	if !exists {
		err = rcc.redisClient.Expire(key, ttl).Err()
	}
	if err != nil {
		log.Errorf("error while setting IncrementRequestCounter : %s for key: %s", err.Error(), key)
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

// FetchCounterValueForKeys fetches the value for counters for given keys
func (rcc *RequestCounterCache) FetchCounterValueForKeys(keys ...string) (counterResponseList []int, err error) {
	log.Infof("fetching values for multiple keys")
	counterValList, err := rcc.redisClient.MGet(keys...).Result()
	log.Infof("counter val is %v", counterValList)
	if err != nil {
		log.Errorf("Error while fetching counter value for key : %s", keys)
		return
	}
	for _, counter := range counterValList {

		val, ok := counter.(string)
		if !ok {
			val = "0"
		}
		intVal, _ := strconv.Atoi(val)
		counterResponseList = append(counterResponseList, intVal)
	}
	return
}
