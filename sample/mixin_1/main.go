package main

import "fmt"

// Father 定义父类
type Father struct {
	Name string
}

func (f *Father) Gamble() {
	fmt.Printf("%s 正在赌博: 梭哈！\n", f.Name)
}

// Son 定义子类，通过匿名嵌入实现 Mixin
type Son struct {
	Father // 匿名嵌入，这就是所谓的 Mixin 组合
	Age    int
}

// Son 也可以有自己的方法
func (s *Son) Study() {
	fmt.Printf("%d 岁的 %s 正在学习 Golang\n", s.Age, s.Name)
}

/**


 */

func main() {
	// 初始化：注意嵌入结构体的赋值方式
	pSon1 := &Son{
		Father: Father{
			Name: "Father",
		},
		Age: 18,
	}

	pSon1.Name = "f Father"

	fmt.Println("")
	fmt.Println("===================開始====================")
	fmt.Println("")

	fmt.Println(pSon1.Name)
	pSon1.Gamble()
	pSon1.Study()

	fmt.Println("")
	fmt.Println("===========================================")
	fmt.Println("")

	pFather2 := &Father{
		Name: "Father",
	}

	pSon2 := &Son{
		Father: *pFather2,
		Age:    18,
	}

	fmt.Println("")
	fmt.Println("===========================================")
	fmt.Println("")
	fmt.Println("pSon2=", pSon2)

	oFather3 := Father{
		Name: "Father",
	}

	pSon3 := &Son{
		Father: oFather3,
		Age:    18,
	}

	fmt.Println("")
	fmt.Println("===========================================")
	fmt.Println("")
	fmt.Println("pSon3=", pSon3)

}
