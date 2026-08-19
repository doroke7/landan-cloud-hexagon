package main

import "fmt"

type User struct {
	Name string
}

func main() {
	sString1 := "AAAAA"
	sString2 := sString1

	fmt.Println(sString1, sString2)

	oUser1 := &User{Name: "!!!!!"}
	oUser2 := oUser1

	oUser1.Name = "BBBBBB"
	fmt.Println(oUser1, oUser2)

	/**
	如果我今天 宣告了一個string變量 string1
	然後我又給他複製給 string2 ， 那在函數編程時間-我就有兩個東西 獨立的 ，各種代表自己的 字串符號連

	但是，如果我今天 宣告了一個 object 變量 oUser1
	然後我又給他拷貝給 oUser2 ，
	在自然界裡面，這兩個 物件 其實都是代表同一個人 同一個自然界的物品，同時修改是比較直觀的。


	*/
}
