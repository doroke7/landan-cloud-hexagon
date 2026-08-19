package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"https://127.0.0.1/App/Resource/AppUser/1",
		"https://127.0.0.1/App/Resource/AppUser/2",
		"https://127.0.0.1/App/Resource/AppUser/3",
		"https://127.0.0.1/App/Resource/AppUser/4",
		"https://127.0.0.1/App/Resource/AppUser/5",
		"https://127.0.0.1/App/Resource/AppUser/6",
		"https://127.0.0.1/App/Resource/AppUser/7",
		"https://127.0.0.1/App/Resource/AppUser/8",
	}

	var wg sync.WaitGroup

	// HTTP client（含 timeout + 忽略 TLS 验证）
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 本地开发用
			},
		},
	}

	for _, url := range urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()

			resp, err := client.Get(u)
			if err != nil {
				fmt.Println("error:", u, err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("read error:", u, err)
				return
			}

			fmt.Printf("URL: %s\nStatus: %s\nBody: %s\n\n",
				u, resp.Status, string(body))
		}(url)
	}

	wg.Wait()
	fmt.Println("All requests done")
}
