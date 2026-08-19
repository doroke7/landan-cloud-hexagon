package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// singleflight 是高併發情況只允許 同時間 一個函數在執行
// Once 不只同時間執行一次，並且永遠只有執行一次

var g singleflight.Group

func fetchUser(userID string) (string, error) {
	// 同一個 key 的同時呼叫，只會執行 fn 一次。
	value, err, shared := g.Do(userID, func() (any, error) {
		fmt.Println("真的去查資料：", userID)

		// 模擬打 DB 或外部 API
		time.Sleep(1 * time.Second)

		return "User data for " + userID, nil
	})

	if err != nil {
		return "", err
	}

	fmt.Printf("userID=%s, shared=%v\n", userID, shared)

	return value.(string), nil
}

func main() {
	var wg sync.WaitGroup

	// 同時 5 個請求查同一個 userID。
	for i := 1; i <= 5; i++ {
		wg.Add(1)

		go func(requestID int) {
			defer wg.Done()

			user, err := fetchUser("user:42")
			if err != nil {
				fmt.Println("error:", err)
				return
			}

			fmt.Printf("request %d 得到：%s\n", requestID, user)
		}(i)
	}

	wg.Wait()
}