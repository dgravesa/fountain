package redis

import (
	"encoding/json"
	"fmt"

	"github.com/dgravesa/fountain/pkg/fountain"
	"github.com/go-redis/redis"
)

// Reservoir is a redis-based implementation of a reservoir
type Reservoir struct {
	clientBase
}

// NewReservoir returns a new Redis-based reservoir
func NewReservoir(addr string) (*Reservoir, error) {
	var err error
	reservoir := new(Reservoir)
	reservoir.clientBase, err = makeClientBase(addr)
	return reservoir, err
}

func uwlKey(userID string) string {
	return fmt.Sprintf("users/%s/wls", userID)
}

func userWls(rdb *redis.Client, userID string) ([]fountain.WaterLog, error) {
	wls := []fountain.WaterLog{}

	wlsStr, err := rdb.Get(uwlKey(userID)).Result()
	if err == nil {
		// some logs already exist, so append to those
		if unmarshalErr := json.Unmarshal([]byte(wlsStr), &wls); unmarshalErr != nil {
			return nil, unmarshalErr
		}
	} else if err != redis.Nil {
		// error other than key not found
		return nil, err
	}

	return wls, nil
}

// WriteWl adds a new user water log to the store
func (r *Reservoir) WriteWl(userID string, wl *fountain.WaterLog) error {
	rdb := redisClient(r)

	// pull water logs for user
	wls, err := userWls(rdb, userID)
	if err != nil {
		return err
	}

	// append this water log and set key-value
	wls = append(wls, *wl)
	wlsBytes, _ := json.Marshal(&wls)
	wlsStr := string(wlsBytes)
	return rdb.Set(uwlKey(userID), wlsStr, 0).Err()
}

// UserWls returns all water logs for a user
func (r *Reservoir) UserWls(userID string) ([]*fountain.WaterLog, error) {
	// pull water logs for user
	wls, err := userWls(redisClient(r), userID)
	if err != nil {
		return nil, err
	}

	// convert to array of pointers
	wlptrs := []*fountain.WaterLog{}
	for i := range wls {
		wlptrs = append(wlptrs, &wls[i])
	}

	return wlptrs, nil
}
