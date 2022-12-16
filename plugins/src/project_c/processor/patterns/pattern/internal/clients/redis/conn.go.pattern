package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func New(opts *redis.Options) (*redis.Client, error) {
	c := redis.NewClient(opts)

	res, err := c.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "error checking connection to redis")
	}

	if res != "pong" {
		return nil, errors.New(fmt.Sprintf("not a pong has returned but %s", res))
	}

	return c, nil
}
