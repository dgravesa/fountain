package redis

import "github.com/go-redis/redis"

type clientBase interface {
	Address() string
}

func makeClientBase(addr string) (clientBase, error) {
	c := clientBaseImpl{address: addr}

	// test redis connection with ping
	rdb := redisClient(c)
	if err := rdb.Ping().Err(); err != nil {
		return nil, err
	}

	return &c, nil
}

func redisClient(c clientBase) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.Address(),
		Password: "",
		DB:       0,
	})
}

type clientBaseImpl struct {
	address string
}

func (c clientBaseImpl) Address() string {
	return c.address
}
