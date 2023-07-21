package lock

import (
	"context"
	"oom_demo/db/redis"
	"sync"
	"testing"
	"time"
)

var (
	rdsClient = redis.GetRedisClient(
		redis.Config{
			Host: "123.57.167.85",
			Port: "6379",
			Pwd:  "1k2_3_456",
		},
	)
	key = "test"
)

func TestMain(m *testing.M) {
	rdsClient.Del(context.Background(), key).Result()
	m.Run()
}
func TestRedisLock(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(i) * 1000 * time.Millisecond)
			defer wg.Done()
			ctx := context.Background()
			redisLock := NewRedisLock(ctx, rdsClient, key)
			id := redisLock.ID()
			if redisLock.Lock() {
				defer redisLock.Unlock()
				t.Log(id, "get lock success")
			} else {
				t.Log(id, "get lock failed")
			}
			time.Sleep(time.Duration((10-i)/2) * 1000 * time.Millisecond)
			t.Log(id, "stop")
		}(i)
	}
	wg.Wait()
}

// 续约测试

func TestRedisLock_KeepAlive(t *testing.T) {
	f := func(sleep int) {
		ctx := context.Background()
		redisLock := NewRedisLock(ctx, rdsClient, key)
		id := redisLock.ID()
		if redisLock.Lock() {
			defer redisLock.Unlock()
			t.Log(id, "get lock success")
		} else {
			t.Log(id, "get lock failed")
		}
		time.Sleep(time.Duration(sleep) * time.Second)
	}
	f2 := func() {
		ctx := context.Background()
		redisLock := NewRedisLock(ctx, rdsClient, key)
		id := redisLock.ID()
		for {
			if redisLock.Lock() {
				defer redisLock.Unlock()
				t.Log(id, "get lock success")
				time.Sleep(1 * time.Second)

				return
			} else {
				t.Log(id, "get lock failed")
			}
			time.Sleep(1 * time.Second)
		}

	}
	go f(10)
	go f2()

	// 20s 后停止
	time.Sleep(20 * time.Second)
}
