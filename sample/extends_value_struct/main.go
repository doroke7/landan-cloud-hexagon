package main

import "fmt"

type Machine3 struct {
	Name string
}

func (oSelf *Machine3) PowerOn() {
	fmt.Println(oSelf.Name + " 我開機了！")
}

type Car3 struct {
	Machine3 // 等價 Machine3: *Machine3
	Color    string
}

type SportsCar3 struct {
	Car3  // 等價 Car3: *Car3
	Speed int
}

func NewMachine3() Machine3 {
	return Machine3{}
}

func NewCar3() Car3 {
	return Car3{}
}

func NewSportsCar3() SportsCar3 {
	return SportsCar3{}
}

func main() {

	oCar3 := NewSportsCar3()
	fmt.Println(oCar3)

}
