package main

import "fmt"

type Animal struct {
	Name string
}

func (oAnimal *Animal) Run() {
	fmt.Printf("%s 正在 run ！\n", oAnimal.Name)
}

type Person struct {
	*Animal
	Age int
}

func (oPerson *Person) Talk() {
	fmt.Printf("%d 岁的 %s 正在 talk\n", oPerson.Age, oPerson.Name)
}

func main() {
	// 初始化：注意嵌入结构体的赋值方式

	pPerson1 := &Person{
		Animal: &Animal{
			Name: "Animal",
		},
		Age: 18,
	}

	fmt.Println(pPerson1.Name)

	fmt.Println("")
	fmt.Println("===========================================")
	fmt.Println("")

	pAnimal2 := &Animal{
		Name: "Animal",
	}
	pPerson2 := &Person{
		Animal: pAnimal2,
		Age:    18,
	}

	fmt.Println(pPerson2.Name)

	fmt.Println("")
	fmt.Println("===========================================")
	fmt.Println("")

	oAnimal2 := Animal{
		Name: "Animal",
	}
	pPerson3 := &Person{
		Animal: &oAnimal2,
		Age:    18,
	}

	fmt.Println(pPerson3.Name)

}
