package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func worker(stop <-chan struct{}) {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {

		select {

		case <-stop:
			fmt.Println("worker shutdown")
			return

		default:

		}

		fmt.Println("sending request")

		resp, err := client.Get(
			"https://httpbin.org/get",
		)

		if err != nil {

			fmt.Println(
				"request error:",
				err,
			)

			time.Sleep(time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)

		resp.Body.Close()

		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		fmt.Println(
			"response:",
			string(body[:50]),
		)

		time.Sleep(time.Second)
	}
}

func main() {

	stop := make(chan struct{})

	go worker(stop)

	// 模擬運行 10 秒
	time.Sleep(10 * time.Second)

	fmt.Println("send shutdown")

	// broadcast shutdown
	close(stop)

	time.Sleep(2 * time.Second)

	fmt.Println("main exit")
}
