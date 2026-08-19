package main

import (
	"fmt"
	"strconv"
)

func main() {
	var input1, input2 string

	// 1. 獲取第一個數字
	fmt.Print("請輸入第一個數字: ")
	fmt.Scanln(&input1)
	num1, err1 := strconv.ParseFloat(input1, 64)
	if err1 != nil {
		fmt.Println("錯誤：第一個輸入不是數字")
		return
	}

	// 2. 獲取第二個數字
	fmt.Print("請輸入第二個數字: ")
	fmt.Scanln(&input2)
	num2, err2 := strconv.ParseFloat(input2, 64)
	if err2 != nil {
		fmt.Println("錯誤：第二個輸入不是數字")
		return
	}

	// 3. 直接相乘並輸出
	result := num1 * num2
	fmt.Printf("\n結果： %v * %v = %v\n", num1, num2, result)
}
