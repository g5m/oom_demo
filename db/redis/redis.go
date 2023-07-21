package redis

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host string
	Port string
	Pwd  string
}
type Client struct {
	*redis.Client
}

var once sync.Once
var client *Client

type redisLogger struct {
	logger *log.Logger
}

func (l *redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}
func NewRedisLogger() *redisLogger {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	return &redisLogger{logger}
}

func GetRedisClient(conf ...Config) *Client {
	once.Do(func() {
		if conf == nil {
			panic("redis config is nil")
		}
		addr := conf[0].Host + ":" + conf[0].Port
		pwd := conf[0].Pwd
		redis.SetLogger(NewRedisLogger())
		c := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pwd, // no password set
			DB:       0,   // use default DB
		},
		)

		if c.Ping(context.Background()).Err() != nil {
			panic("redis connect error")
		}

		client = &Client{c}
	})
	return client
}
