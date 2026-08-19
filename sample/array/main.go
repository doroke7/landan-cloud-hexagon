package main

import "fmt"

// Golang中，函數的修改 不會 改變 array 外部的值，因為他是 call by value
//  array 定義的長度 後來絕對不能擴張

func UpdateFirstElementOfArray(sNewName string, aData [3]string) [3]string {
	// 1. 在副本上进行操作
	aData[0] = sNewName

	// 2. 将修改后的副本整体返回
	return aData
}

func main() {
	// 定义原始数组
	aTags1 := [3]string{"Go", "PHP", "Java"}

	// 调用函数。注意：必须用返回值重新赋值，否则修改就丢失了
	aTags2 := UpdateFirstElementOfArray("Python", aTags1)
	fmt.Println("1. 原始数据:", aTags1)

	fmt.Println("2. 操作后果:", aTags2)
}
