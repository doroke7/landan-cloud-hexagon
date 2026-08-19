package main

import "fmt"

func main() {
	// 1. 初始化一个 slice
	sTags := []string{"Go", "PHP"}

	// 在 golang 的世界 append 一定要返回
	sTags = append(sTags, "Java")
	sTags = append(sTags, "Python", "C++")

	fmt.Println(sTags) // [Go PHP Java Python C++]
}
