package redis

import (
	"context"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host string
	Port int
}
type Client struct {
	clients *redis.ClusterClient
}

var once sync.Once
var client *Client

func GetRedisClient(conf ...Config) *Client {
	once.Do(func() {
		if conf == nil {
			panic("redis config is nil")
		}
		addrs := make([]string, 0)
		for _, v := range conf {
			addrs = append(addrs, v.Host+":"+strconv.Itoa(v.Port))
		}
		c := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: addrs,
		})
		if c.Ping(context.Background()).Err() != nil {
			panic("redis connect error")
		}
		client = &Client{clients: c}
	})
	return client
}


