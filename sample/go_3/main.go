package main

import (
	"fmt"
	"time"
)

func DoSomething() bool {
	// 啟動背景小弟
	ch := make(chan bool)
	go func() {
		time.Sleep(5 * time.Second) // 模擬偷懶兩秒
		ch <- true
		fmt.Println("【背景】小弟默默把事情做完了！")
	}()
	select {
	case bResult := <-ch:
		return bResult
	case <-time.After(1 * time.Second): // 5 秒 goruntine 才丟 channel ，但是 1秒已經return了
		return false
	}

}

// 造成 goruntine 一直釋放不了

func main() {
	// goruntine 是會脫離父親 func 獨自運行的
	result := DoSomething()
	fmt.Println(result) // 會先印出主函數下班

	time.Sleep(3 * time.Second) // 留時間讓背景小弟把話說完
}
