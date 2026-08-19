package main

import "fmt"

func main() {
	var i interface{} = "Hello Go" // 盒子裡裝了一個字串

	// 語法：值, 是否成功 := 接口變數.(目標型別)
	s, ok := i.(string)

	if ok {
		fmt.Println("轉換成功:", s) // 輸出: Hello Go
	} else {
		fmt.Println("轉換失敗：這不是字串")
	}
}
