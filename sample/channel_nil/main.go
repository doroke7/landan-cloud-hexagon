package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 42
		close(ch1)
	}()

	go func() {
		ch2 <- 99
		close(ch2)
	}()

	// 我們想要同時監聽 ch1 與 ch2，直到兩個都關閉為止
	for ch1 != nil || ch2 != nil {
		select {
		case v, ok := <-ch1:
			if !ok {
				// ch1 已經關閉，將它設為 nil！
				// 這樣下一次 loop 時，select 就會自動忽略這個 case（因為讀取 nil 會阻塞）
				ch1 = nil
				fmt.Println("ch1 已關閉，停用 ch1 監聽")
				continue
			}
			fmt.Println("收到 ch1:", v)

		case v, ok := <-ch2:
			if !ok {
				// ch2 已經關閉，將它設為 nil
				ch2 = nil
				fmt.Println("ch2 已關閉，停用 ch2 監聽")
				continue
			}
			fmt.Println("收到 ch2:", v)
		}
	}

	fmt.Println("所有 Channel 處理完畢！")
}
