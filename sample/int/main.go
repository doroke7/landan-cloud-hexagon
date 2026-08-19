package main

import "fmt"

// int 是一種跟著 OS 的型態
// 32 位元 os 他幾乎等於 int32
// 64 位元 os 他幾乎等於 int64
func main() {
	var iA int = 10
	var iB int64 = 10

	// iC := iA + iB 即使在 golang 64 電腦， int64 int 也不能直接相加

	iC := int64(iA) + iB // ✅ 類型轉換後可運算

	fmt.Printf("iC: %d", iC)
}
