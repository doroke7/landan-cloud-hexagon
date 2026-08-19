package main

import (
	"fmt"
	"time"
)

type Request struct {
	ID              int
	Value           int
	ResponseChannel chan Response
}

type Response struct {
	ID     int
	Result int
}

/*
   這種設計模式的實作是：
   client 生成 N 個 request元素的 channel 給 出去，
   client 並且在 1個 request元素 綁定 一個 response channel

   server 批量處理這個 N 個 request 並且 把 response channel 資料塞給 request 的屬性

   client 再去接 這個 response 異步處理

*/

func server(oRequest Request) {

	fmt.Println("server processing request:", oRequest.ID)

	// 模擬工作
	time.Sleep(500 * time.Millisecond)

	// 回傳結果
	oRequest.ResponseChannel <- Response{
		ID:     oRequest.ID,
		Result: oRequest.Value * 2,
	}

	close(oRequest.ResponseChannel)
}

func client(id int, oRequestChannel chan<- Request) {

	oResponsetChannel := make(chan Response)

	// 發送 request
	oRequestChannel <- Request{
		ID:              id,
		Value:           id * 10,
		ResponseChannel: oResponsetChannel,
	}

	// 等待 response 或 timeout
	select {

	case result := <-oResponsetChannel:
		fmt.Println(
			"client",
			id,
			"got result:",
			result.Result,
		)

	case <-time.After(time.Second):
		fmt.Println(
			"client",
			id,
			"timeout",
		)
	}
}

func main() {

	oRequestChannel := make(chan Request)

	// 啟動 server

	// 多個 client
	for i := 1; i <= 5; i++ {
		go client(i, oRequestChannel)
	}

	for oRequest := range oRequestChannel {
		go server(oRequest)
	}

	// 等待 goroutine
	time.Sleep(5 * time.Second)

	close(oRequestChannel)
}
