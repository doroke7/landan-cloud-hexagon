package main

import "fmt"

type User struct {
	Name string
}

func main() {
	// 第一種：宣告為「值」 (Value)
	var u1 User
	u1.Name = "小明"

	// 第二種：宣告為「指標」 (Pointer)
	u2 := &User{Name: "小美"}

	fmt.Println(u1, u2)

	// u1 的宣告，如果未指定，會給zero-value
	// u2 的宣告，如果未指定，會給 nil pointer
}
