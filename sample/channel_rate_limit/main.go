package main

import (
	"fmt"
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(rate int, capacity int) *RateLimiter {

	rl := &RateLimiter{
		tokens: make(chan struct{}, capacity),
	}

	// 初始 token
	for i := 0; i < capacity; i++ {
		rl.tokens <- struct{}{}
	}

	// 補充 token
	go func() {

		ticker := time.NewTicker(
			time.Second / time.Duration(rate),
		)

		defer ticker.Stop()

		for range ticker.C {

			select {

			case rl.tokens <- struct{}{}:
				// 補 token

			default:
				// token 滿了
			}
		}

	}()

	return rl
}

func (oSelf *RateLimiter) Limit() {

	// 沒 token 就阻塞
	<-oSelf.tokens
}

func main() {

	// 每秒產生 5 個 token
	// 最大容量 3
	limiter := NewRateLimiter(5, 3)

	for i := 0; i < 10; i++ {

		go func(id int) {

			limiter.Limit()

			fmt.Println(
				"request",
				id,
				time.Now().Format("15:04:05"),
			)

		}(i)
	}

	time.Sleep(3 * time.Second)
}
