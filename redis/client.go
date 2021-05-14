package redis

import (
	"github.com/go-redis/redis"
)

type Client struct{ cli *redis.Client }

var Cli *Client

func InitRedis() error {
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	Cli = &Client{r}
	return r.Ping().Err()
}

func (r Client) Get(key string) (res []byte, ok bool) {
	res, err := r.cli.Get(key).Bytes()
	ok = err == nil
	return
}

func (r Client) HGet(key string, hash string) (res []byte, ok bool) {
	res, err := r.cli.HGet(key, hash).Bytes()
	ok = err == nil
	return
}
