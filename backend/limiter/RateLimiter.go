package limiter

import (
	"context"
	"time"
)

type RateLimiter struct {
	TokenBucket    chan struct{}
	Context        context.Context
	CancelFunction func()
}

func NewRateLimiter(ratePerSecond int) *RateLimiter {
	ctx, cancelFunc := context.WithCancel(context.Background())

	limiter := &RateLimiter{
		TokenBucket:    make(chan struct{}, ratePerSecond),
		Context:        ctx,
		CancelFunction: cancelFunc,
	}

	go func() {
		for {
			time.Sleep(time.Second / time.Duration(ratePerSecond))
			select {
			case limiter.TokenBucket <- struct{}{}:
			case <-limiter.Context.Done():
				return
			default:
			}
		}
	}()

	return limiter
}

func (rl *RateLimiter) Destroy() {
	rl.CancelFunction()
}
