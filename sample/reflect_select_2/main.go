package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 建立兩個測試用的 Channel
	ch1 := make(chan int, 1)
	ch2 := make(chan string, 1)

	// 塞入測試資料
	ch1 <- 42
	ch2 <- "Hello, Go!"
	ch1 <- 42
	ch1 <- 43
	ch1 <- 55
	ch1 <- 42
	ch1 <- 77

	// 1. 建立一個 reflect.SelectCase 的 Slice
	cases := make([]reflect.SelectCase, 3) // 2 個 Channel + 1 個 Default

	// 2. 設定第一個 Case: 接收 ch1
	cases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ch1),
	}

	// 3. 設定第二個 Case: 接收 ch2
	cases[1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ch2),
	}

	// 4. 設定第三個 Case: Default (非阻塞)
	cases[2] = reflect.SelectCase{
		Dir: reflect.SelectDefault,
	}

	// 5. 執行 reflect.Select
	// chosen: 被選中的 case 索引
	// recv: 接收到的值 (如果是 Recv 操作)
	// recvOK: 接收是否成功 (Channel 是否未關閉)
	chosen, recv, recvOK := reflect.Select(cases)

	switch chosen {
	case 0:
		if recvOK {
			fmt.Printf("從 ch1 收到資料: %v (型態: %s)\n", recv.Int(), recv.Type())
		} else {
			fmt.Println("ch1 已關閉")
		}
	case 1:
		if recvOK {
			fmt.Printf("從 ch2 收到資料: %v (型態: %s)\n", recv.String(), recv.Type())
		} else {
			fmt.Println("ch2 已關閉")
		}
	case 2:
		fmt.Println("執行了 Default 分支（沒有 Channel 準備好）")
	}


	// 
}