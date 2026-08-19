package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*

  go 喜歡用最小組件 context
  來配置 給所有 tcp/udp 連線 timeout，
  而不是 在 http 上配置 timeout 屬性。
*/

func main() {
	// 建立一個 20 秒 Timeout 的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 建立 HTTP Request，並將 Context 綁定到 Request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://example.com",
		nil,
	)
	if err != nil {
		panic(err)
	}

	// 發送 HTTP Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("request failed:", err)
		return
	}
	defer resp.Body.Close()

	// 讀取 Response Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body failed:", err)
		return
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println(string(body))
}
