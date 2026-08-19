package main

import "fmt"

/*
   基本上 函數下面的 aNumbers 型態一樣
   函數外面的用法不同 ，一個是 逗號隔開， 一個是塞array
*/

func Sum1(aNumbers []int) int {
	iSum := 0
	for _, iNumber := range aNumbers {
		iSum += iNumber
	}
	return iSum
}

func Sum2(aNumbers ...int) int {
	iSum := 0
	for _, iNumber := range aNumbers {
		iSum += iNumber
	}
	return iSum
}

func main() {
	aNumbers := []int{1, 2, 3, 4, 5}

	fmt.Println("Sum1(aNumbers):", Sum1(aNumbers))

	fmt.Println("Sum2(1, 2, 3):", Sum2(1, 2, 3))
	fmt.Println("Sum2(aNumbers...):", Sum2(aNumbers...))
}
