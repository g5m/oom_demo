package rate_limit

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"golang.org/x/time/rate"
)

func Test_tokenBucket_startSetToken(t *testing.T) {
	type fields struct {
		Capacity        int32
		AvailableTokens int32
		Rate            int32
		name            string
		stop            chan struct{}
		m               sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &tokenBucket{
				Capacity:        tt.fields.Capacity,
				AvailableTokens: tt.fields.AvailableTokens,
				Rate:            tt.fields.Rate,
				name:            tt.fields.name,
				stop:            tt.fields.stop,
				m:               tt.fields.m,
			}
			tb.startSetToken()
		})
	}
}

func Test_tokenBucket(t *testing.T) {
	tokenBucket := NewTokenBucket("test", 100, 100)
	g := sync.WaitGroup{}
	g.Add(200)

	for i := 0; i < 200; i++ {
		go func(i int) {
			defer g.Done()
			if tokenBucket.Take() {
				fmt.Printf("goroutine %d get token\n", i)
			} else {
				fmt.Printf("goroutine %d get token failed\n", i)
			}
		}(i)
	}
	g.Wait()
}

func TestGoLimiter(t *testing.T) {
	fmt.Printf("start\n")
	g := sync.WaitGroup{}
	g.Add(200)
	limiter := rate.NewLimiter(10, 50)

	for i := 0; i < 200; i++ {
		go func(i int) {
			defer g.Done()
			limiter.Wait(context.Background())
			fmt.Printf("goroutine %d get token\n", i)
			//limiter.Wait(context.Background())
			// if limiter.Allow() {
			// 	time.Sleep(100 * time.Microsecond)
			// 	fmt.Printf("goroutine %d get token\n", i)
			// 	//log.Println("goroutine %d get token\n", i)
			// } else {
			// 	//log.Println("goroutine %d get token failed\n", i)
			// 	fmt.Printf("goroutine %d get token failed\n", i)
			// }

		}(i)
	}
	g.Wait()
	// limiter.Allow()
}
