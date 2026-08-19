package main

import "fmt"

func test() int {
	iI := 0

	// defer 的執行是 func return 之後 按照倒序執行
	defer func() {
		iI += 1
		fmt.Println("defer-a", iI)
	}()
	defer func() {
		iI += 1
		fmt.Println("defer-b", iI)
	}()

	defer func() {
		iI += 1
		fmt.Println("defer-c", iI)
	}()

	defer func() {
		iI += 1
		fmt.Println("defer-d", iI)
	}()

	defer func() {
		iI += 1
		fmt.Println("defer-e", iI)
	}()
	return iI
}

func main() {
	test()
}
