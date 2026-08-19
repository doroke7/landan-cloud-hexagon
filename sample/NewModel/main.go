package main

import "fmt"

type Model struct {
	Name string
	Age  int
}

func NewModel1(name string) *Model {
	// NewXXYY （&XXYY{}） 並不是單例， 他是每次都聲請一個 放 XXYY 新實例的內存空間
	// 差別在於 2次拷貝 “值” 的時候， 指標的會一起修改

	return &Model{
		Name: name,
	}
}

func NewModel2(name string) Model {
	return Model{
		Name: name,
	}
}

func main() {
	u1 := NewModel1("AAAA")
	u2 := u1

	u1.Name = "Gemini"

	fmt.Println("u1=", u1) //
	fmt.Println("u2=", u2) //  u2 u1 會一起修改

	m1 := NewModel2("AAAA")
	m2 := m1
	m1.Name = "CAR"

	fmt.Println("m1=", m1) //
	fmt.Println("m2=", m2) //  m2 m1 不會一起修改
}
