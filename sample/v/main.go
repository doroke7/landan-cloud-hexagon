package main

import "fmt"

func NewUser(sName string) *User {
	return &User{
		Name: sName,
	}
}

type User struct {
	Name string
}

func main() {

	oUser := NewUser("Tom")

	fmt.Printf("%v", oUser) // go
	fmt.Println("")

	fmt.Printf("%+v", oUser) // go 裡面 + 代表 more
	fmt.Println("")

	fmt.Printf("%#v", oUser) // go 裡面 # 代表 detail
	fmt.Println("")

	fmt.Printf("%T", oUser)
	fmt.Println("")

	fmt.Printf("%t", oUser)
	fmt.Println("")
}
