package main

import (
	"fmt"
	"time"
)

func DoSomething() string {
	// 啟動背景小弟
	go func() {
		time.Sleep(1 * time.Second) // 模擬偷懶兩秒
		fmt.Println("【背景】小弟默默把事情做完了！")
	}()

	// 主函數直接拍拍屁股先下班 return
	return "【主函數】我先下班囉！"
}

func main() {
	// goruntine 是會脫離父親 func 獨自運行的
	result := DoSomething()
	fmt.Println(result) // 會先印出主函數下班

	time.Sleep(3 * time.Second) // 留時間讓背景小弟把話說完
}
