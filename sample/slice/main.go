package main

import "fmt"

// go 的 Slice 本質是帶有指針的 array
// 如果函數操作了 slice 的元素，會改變原來的值
// 如果函數 append 就會產生一個新的（這個是append 的特性）

func UpdateFirstElementOfSlice(aRows []string) {
	// 1. 修改元素：这会直接改变原始底层数组
	aRows[0] = "Python"

	// 2. 这里的修改在 main 函数里是可见的
	fmt.Println("函数内部修改完毕")
}

func AppendForSlice(aRows []string, sString string) {
	// 1. 修改元素：这会直接改变原始底层数组
	aRows = append(aRows, sString)

	// 2. 这里的修改在 main 函数里是可见的
	fmt.Println("函数内部修改完毕")
}

func main() {
	// 定义一个 Slice
	aTags := []string{"Go", "PHP", "Java"}
	aNames := []string{"Tom", "Mary", "Jim"}

	fmt.Println("1. 原始数据 aTags:", aTags)

	// 直接传 slice，不需要 &
	UpdateFirstElementOfSlice(aTags)

	// 最终结果：你会发现 sTags[0] 已经变成 Python 了
	fmt.Println("2. 操作后果 aTags:", aTags)

	AppendForSlice(aNames, "Laplace")

	fmt.Println("2. 操作后果 aNames :", aNames)

}
