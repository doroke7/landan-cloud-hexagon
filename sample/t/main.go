package main

import "fmt"

func main() {
	// 編譯時候判斷：利用 %T 看看 型態
	var a = 10
	var b = "hello"

	fmt.Printf("a type: %T\n", a)
	fmt.Printf("b type: %T\n", b)
}
