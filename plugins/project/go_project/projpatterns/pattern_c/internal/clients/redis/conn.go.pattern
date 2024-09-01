package redis

import (
	"strconv"

	"github.com/Red-Sock/toolbox/closer"
	"github.com/Red-Sock/trace-errors"
	"github.com/go-redis/redis"
	"github.com/godverv/matreshka/resources"
)

var ErrUnexpectedPing = errors.New("error pinging redis")

type Client redis.Client

func New(cfg *resources.Redis) (*Client, error) {
	opts := &redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(int(cfg.Port)),
		Password: cfg.Pwd,
		DB:       cfg.Db,
	}
	c := redis.NewClient(opts)

	res, err := c.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "error checking connection to redis")
	}

	if res != "pong" {
		return nil, errors.Wrapf(ErrUnexpectedPing, "not a pong has returned but %s", res)
	}

	closer.Add(c.Close)

	return (*Client)(c), nil
}
