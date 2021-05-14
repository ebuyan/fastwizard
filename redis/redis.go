package redis

import (
	"github.com/go-redis/redis"
)

type Redis struct{ cli *redis.Client }

func NewRedis() Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return Redis{rdb}
}

func (r Redis) Get(key string) (res []byte, ok bool) {
	res, err := r.cli.Get(key).Bytes()
	ok = err == nil
	return
}

func (r Redis) HGet(key string, hash string) (res []byte, ok bool) {
	res, err := r.cli.HGet(key, hash).Bytes()
	ok = err == nil
	return
}
