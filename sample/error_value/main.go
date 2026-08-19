package main

import (
	"errors"
	"fmt"
)

// 你這樣說 error value 是為了 提升程序基本性能, try catch 開銷比較大

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		// 方案 A：回傳一個簡單的錯誤對象
		return 0, errors.New("除數不能為零")
	}
	return a / b, nil
}

func main() {
	result, err := Divide(10, 0)
	if err != nil {
		fmt.Println("出錯了:", err)
		return
	}
	fmt.Println("結果:", result)
}
