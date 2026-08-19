package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		go func(i int) {
			fmt.Println(i) //
		}(i)
	}

	time.Sleep(time.Second * 10) // 停頓，讓 goroutine 繼續跑
}
