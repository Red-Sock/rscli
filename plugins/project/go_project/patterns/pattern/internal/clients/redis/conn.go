package redis

import (
	"strconv"

	"github.com/Red-Sock/toolbox/closer"
	"github.com/go-redis/redis"
	"github.com/godverv/matreshka/resources"
	"go.redsock.ru/rerrors"
)

var ErrUnexpectedPing = rerrors.New("error pinging redis")

func New(cfg *resources.Redis) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(int(cfg.Port)),
		Password: cfg.Pwd,
		DB:       cfg.Db,
	}
	c := redis.NewClient(opts)

	res, err := c.Ping().Result()
	if err != nil {
		return nil, rerrors.Wrap(err, "error checking connection to redis")
	}

	if res != "pong" {
		return nil, rerrors.Wrapf(ErrUnexpectedPing, "not a pong has returned but %s", res)
	}

	closer.Add(c.Close)

	return c, nil
}
