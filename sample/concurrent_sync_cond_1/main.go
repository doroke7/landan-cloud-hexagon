package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

// Signal 通知（唤醒）一个正在 Wait() 的 goroutine
// BroadCast 通知（唤醒）全部正在 Wait() 的 goroutine，但是實際情況有併發限制，全部喚醒後可能大部分還是又wait了

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	done := false
	var result string

	// HTTP goroutine
	go func() {
		resp, err := http.Get("https://httpbin.org/get")
		if err != nil {
			result = err.Error()
		} else {
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			result = string(body)
		}

		// 请求完成后通知
		mu.Lock()
		done = true
		cond.Signal()
		mu.Unlock()
	}()

	// 等待 HTTP 请求完成
	mu.Lock()
	for !done {
		cond.Wait()
	}
	mu.Unlock()

	// 只有请求完成后才会执行
	fmt.Println("HTTP request finished")
	fmt.Println(result)
}
