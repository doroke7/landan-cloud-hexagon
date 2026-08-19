package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

/*
   sync.cond 本質就是有條件的控制流程的程序
*/

func main() {
	aUrls := []string{
		"https://127.0.0.1/App/Resource/AppUser/1",
		"https://127.0.0.1/App/Resource/AppUser/2",
		"https://127.0.0.1/App/Resource/AppUser/3",
		"https://127.0.0.1/App/Resource/AppUser/4",
		"https://127.0.0.1/App/Resource/AppUser/5",
		"https://127.0.0.1/App/Resource/AppUser/6",
		"https://127.0.0.1/App/Resource/AppUser/7",
		"https://127.0.0.1/App/Resource/AppUser/8",
	}

	var mu sync.Mutex
	oCond := sync.NewCond(&mu)

	nMaxConcurrent := 2 // 并发上限（可调）
	nRunning := 0       // 当前运行中的 goroutine 数
	nCompleted := 0     // 已完成数量
	nUrlLen := len(aUrls)

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 本地测试用
			},
		},
	}

	for _, sUrl := range aUrls {
		mu.Lock()
		for nRunning >= nMaxConcurrent {
			oCond.Wait() // 等待有空位
		}
		nRunning++
		mu.Unlock()

		go func(sU string) {
			// 执行请求
			resp, err := client.Get(sU)
			if err != nil {
				fmt.Println("error:", sU, err)
			} else {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				fmt.Printf("URL: %s\nStatus: %s\nBody: %s\n\n", sU, resp.Status, string(body))
			}

			// 完成后更新状态
			mu.Lock()
			nRunning--
			nCompleted++
			oCond.Signal() // 通知一个等待的 goroutine

			// 如果全部完成，唤醒主线程
			if nCompleted == nUrlLen {
				oCond.Broadcast()
			}
			mu.Unlock()
		}(sUrl)
	}

	// 主线程等待所有任务完成
	mu.Lock()
	for nCompleted < nUrlLen {
		oCond.Wait()
	}
	mu.Unlock()

	fmt.Println("All requests done")
}
