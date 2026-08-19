package main

import (
	"fmt"
)

func main() {

	/*
		為什麼 any/interface{} 1換 “1” 不會報錯？
		因為任何型別都滿足：空介面的 0 個方法

	*/

	var iN interface{} = 5

	fmt.Println(iN)

	iN = "1111"

	fmt.Println(iN)

}
