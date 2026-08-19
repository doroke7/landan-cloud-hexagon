package main

import "fmt"

func main() {
	iNumber := 1
	sString := "22"

	check(iNumber)
	check(sString)
}

func check(i interface{}) {

	/**
	簡單的說 .type 是一種golang 獨有的特殊 語法，他只能用在 switch case
	在 case 判斷裡面 他是 “型態”，在 case 之後的程序式 他是 “值”
	*/
	switch v := i.(type) {
	case int:
		// 在這裡 v 就是 int
		fmt.Println("這是數字:", v+1)

	case string:
		// 在這裡 v 就是 string
		fmt.Println("這是字串:", v+"!!")

	case bool:
		// 在這裡 v 就是 bool
		fmt.Println("這是布林:", v)

	default:
		fmt.Println("我不認識這個型態")
	}
}
