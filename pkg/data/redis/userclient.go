package redis

import (
	"encoding/json"
	"fmt"

	"github.com/dgravesa/fountain/pkg/fountain"
	"github.com/go-redis/redis"
)

// UserStore is a client to a Redis-based user store
type UserStore struct {
	address string
}

// NewUserStore returns a new Redis-based client
func NewUserStore(addr string) (*UserStore, error) {
	client := UserStore{address: addr}

	// test redis connection with ping
	rdb := redisClient(&client)
	if err := rdb.Ping().Err(); err != nil {
		return nil, err
	}

	return &client, nil
}

func uKey(uid string) string {
	return fmt.Sprintf("users/%s", uid)
}

func redisClient(c *UserStore) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.address,
		Password: "",
		DB:       0,
	})
}

// PutUser adds a new user to the store
func (c *UserStore) PutUser(user *fountain.User) error {
	rdb := redisClient(c)
	userBytes, _ := json.Marshal(user)
	return rdb.Set(uKey(user.ID), string(userBytes), 0).Err()
}

// User retrieves a user from the store
func (c *UserStore) User(userID string) (*fountain.User, error) {
	rdb := redisClient(c)

	// get bytes from redis
	userStr, err := rdb.Get(uKey(userID)).Result()
	if err == redis.Nil {
		// key not found
		return nil, nil
	} else if err != nil {
		// other error
		return nil, err
	}

	// unmarshal into user struct
	user := new(fountain.User)
	err = json.Unmarshal([]byte(userStr), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
