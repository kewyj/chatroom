package src

import (
	"testing"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(2)
	for {
		select {
		case <-rl.TokenBucket:

		default:

		}
	}
}
