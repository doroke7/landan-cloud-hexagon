package main

import "fmt"

func main() {
	fmt.Println("1. 开启程序 (Try)")

	defer func() {
		// --- 这里相当于 Catch ---
		// 在 defer 中 使用了 recover 后就会 捕获 panic 的 讯息对吧
		// 只要 recover() 在 defer 函数中被调用，它就像一个**“安全气囊”**。当程序发生 panic 时，气囊弹出，程序不会直接崩溃（Crash），而是会把 panic 传递的那条讯息给“接住”。
		if r := recover(); r != nil {
			fmt.Printf("2. 捕获到异常 (Catch): %v\n", r)
		}

		// --- 这里相当于 Finally ---
		fmt.Println("3. 无论如何都会执行 (Finally)")
	}()

	// 模拟一个严重错误：比如数组越界或空指针
	fmt.Println("--- 准备触发 Panic ---")
	panic("panic 手动触发！")

	fmt.Println("这行永远不会被执行")
}
