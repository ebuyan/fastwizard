package redis

import (
	"github.com/go-redis/redis"
)

type Client struct{ cli *redis.Client }

var Cli *Client

func InitRedis() (err error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	err = client.Ping().Err()
	Cli = &Client{client}
	return
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

func (r Client) Close() {
	r.cli.Close()
}
