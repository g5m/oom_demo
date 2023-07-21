package lock

import (
	"context"
	"crypto/rand"
	"log"
	"oom_demo/db/redis"
	"time"
)

// 1. 判断key是否存在，如果不存在，设置key
// 2. 如果存在，判断是否是自己的锁，如果是，返回 1，如果不是，返回0
// 3. 设置key的过期时间并 返回1
var luaScript = `
if redis.call("exists", KEYS[1]) == 0 then
    redis.call("set", KEYS[1], ARGV[1])
else
    if redis.call("get", KEYS[1]) == ARGV[1] then
    else
        return 0
    end
end

redis.call("expire", KEYS[1], ARGV[2])
return 1`

// 分布式锁，需要进行续约

type RedisLock struct {
	id     uint64
	rc     *redis.Client
	expire int           // 过期时间 单位 ms
	close  chan struct{} // 关闭续约
	ctx    context.Context
	key    string
}

// NewRedisLock 初始化redis锁
func NewRedisLock(ctx context.Context, rc *redis.Client, key string) *RedisLock {
	// 获得协程id
	n, _ := rand.Prime(rand.Reader, 64)
	id := n.Uint64()
	return &RedisLock{id: id, rc: rc, expire: 1000, close: make(chan struct{}), ctx: ctx, key: key}
}

// getID
func (r *RedisLock) ID() uint64 {
	return r.id
}

// Lock 加锁
// key: 锁的key
// value: 锁的value
func (r *RedisLock) Lock() bool {
	res, err := r.rc.Eval(r.ctx, luaScript, []string{r.key}, r.id, r.expire).Result()
	if err != nil {
		return false
	}
	if res.(int64) == 1 {
		// 续约
		go r.keepAlive(r.ctx, r.key)
		return true
	}
	return false
}

// Unlock 解锁
func (r *RedisLock) Unlock() {
	r.close <- struct{}{}
	ctx := r.ctx
	if r.ctx.Err() != nil {
		ctx = context.Background()
	}

	r.rc.Del(ctx, r.key)
	log.Printf("%d 解锁", r.id)
}

// keepAlive 续约
func (r *RedisLock) keepAlive(ctx context.Context, key string) {
	ticker := time.NewTicker(time.Duration(r.expire/2) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("%d 续约", r.id)
			r.rc.Expire(ctx, key, time.Duration(r.expire)*time.Millisecond)
		case <-r.close:
			return
		case <-ctx.Done():
			return
		}
	}
}
