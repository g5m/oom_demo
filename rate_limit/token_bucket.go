package rate_limit

import (
	"sync"
	"sync/atomic"
	"time"
	"golang.org/x/time/rate"
)

type tokenBucket struct {
	// 令牌桶的容量
	Capacity int32
	// 令牌桶的剩余数量
	AvailableTokens int32
	// 令牌桶的填充速率
	Rate int32
	name string
	stop chan struct{}
	m    sync.RWMutex
}

var tokenBucketMap sync.Map

// NewTokenBucket 初始化令牌桶
func NewTokenBucket(name string, capacity int32, rate int32) *tokenBucket {
	tb, ok := tokenBucketMap.Load(name)
	if ok {
		return tb.(*tokenBucket)
	}
	newTb := &tokenBucket{
		Capacity:        capacity,
		AvailableTokens: capacity,
		Rate:            rate,
		name:            name,
	}
	tokenBucketMap.Store(name, newTb)
	go newTb.startSetToken()
	return newTb
}

func GetTokenBucket(name string) *tokenBucket {
	tb, ok := tokenBucketMap.Load(name)
	if !ok {
		return nil
	}
	return tb.(*tokenBucket)
}

func (tb *tokenBucket) startSetToken() {
	for {
		select {
		case <-time.After(time.Second):
			tb.set()
		case <-tb.stop:
			return
		}
	}
}

func (tb *tokenBucket) Stop() {
	tokenBucketMap.Delete(tb.name)
	tb.stop <- struct{}{}
}

// Take 从令牌桶中取出一个令牌
func (tb *tokenBucket) Take() bool {
	if tb.AvailableTokens == 0 {
		return false
	}
	atomic.AddInt32(&tb.AvailableTokens, -1)
	return true
}

func (tb *tokenBucket) set() {
	tb.m.Lock()
	defer tb.m.Unlock()
	target := tb.AvailableTokens + tb.Rate
	if target > tb.Capacity {
		tb.AvailableTokens = tb.Capacity
		return
	}
	tb.AvailableTokens = target
}



func Limit() {
	limiter:=rate.NewLimiter(100, 100)
	limiter.Allow()
}