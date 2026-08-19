package main

import "fmt"

func main() {
	// 1. Slice (切片/動態陣列)
	// 建立一個可以存 3 個整數的陣列
	s := make([]int, 3)
	s[0] = 10
	fmt.Println("Slice:", s) // [10 0 0]

	// 2. Map (鍵值對/字典)
	// 建立一個 Key 是字串，Value 是整數的容器
	m := make(map[string]int)
	m["age"] = 30
	fmt.Println("Map:", m) // map[age:30]

	// 3. Channel (通道 - 用於協程溝通)
	// 建立一個可以放 1 個字串的管道
	c := make(chan string, 1)
	c <- "hello"
	fmt.Println("Channel:", <-c) // hello
}

/**
Go 拋棄了傳統的「類別（Class）」與「建構子（Constructor）」，
但依然需要一種機制來確保複雜數據結構在記憶體中是「準備就緒」的狀態。
所以 他用了 make 陣列,
           make Map,
        make chan 的 達到類似效果
*/
