package main

import "fmt"

type Machine2 struct {
	Name string
}

func (oSelf *Machine2) PowerOn() {
	fmt.Println(oSelf.Name + " 我開機了！")
}

type Car2 struct {
	*Machine2 // 等價 Machine2: *Machine2
	Color     string
}

type SportsCar2 struct {
	*Car2 // 等價 Car2: *Car2
	Speed int
}

// 使用 指針做嵌入繼承的時候，需要在 New函數上把父親結構體也 new 上， 不然會造成找不到
// 使用 值做嵌入繼承的時候，則不需要 做 【New函數上把父親結構體也 new 上】，他依然找的到

// 所以，處理 pointer 嵌入的繼承結構體， 除了 在子結構體 嵌入父 結構體 之外， 也需要在 子的new 函數 加上 父的new 函數
// 因為 value的都會有基本函數，pointer 則可能為 nil

func NewMachine2() *Machine2 {
	return &Machine2{}
}

func NewCar2() *Car2 {
	return &Car2{
		Machine2: NewMachine2(),
	}
}

func NewSportsCar2() *SportsCar2 {
	return &SportsCar2{
		Car2: NewCar2(),
	}
}

func main() {

	// 有使用 NewXX 做嵌入繼承的話， 宣告一個子類就會非常簡單

	oCar2 := NewSportsCar2()
	fmt.Println(oCar2)

}
