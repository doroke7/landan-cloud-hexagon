package main

import "fmt"

type Machine1 struct {
	Name string
}

func (oSelf *Machine1) PowerOn() {
	fmt.Println(oSelf.Name + " 我開機了！")
}

type Car1 struct {
	*Machine1 // 等價 Machine: *Machine
	Color     string
}

type SportsCar1 struct {
	*Car1 // 等價 Car: *Car
	Speed int
}

func main() {

	// 不使用 NewXX 做嵌入繼承的話， 宣告一個子類就會非常麻煩

	oCar := &SportsCar1{
		Car1: &Car1{
			Machine1: &Machine1{
				Name: "紅色跑車",
			},
			Color: "red",
		},
		Speed: 10,
	}
	fmt.Println(oCar)

}
