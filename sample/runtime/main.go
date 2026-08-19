package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

func main() {
	// 顯示 CPU 架構
	fmt.Println("runtime.GOARCH(CPU 架構)=", runtime.GOARCH)

	// 判斷系統位元
	if unsafe.Sizeof(uintptr(0)) == 8 {
		fmt.Println("系統位元: 64 位元")
	} else if unsafe.Sizeof(uintptr(0)) == 4 {
		fmt.Println("系統位元: 32 位元")
	} else {
		fmt.Println("未知系統位元")
	}

	var x int
	if unsafe.Sizeof(x) == 8 {
		fmt.Println("int size = 8 bytes → 系統是 64 位元")
	} else if unsafe.Sizeof(x) == 4 {
		fmt.Println("int size = 4 bytes → 系統是 32 位元")
	} else {
		fmt.Println("未知系統位元")
	}
}
