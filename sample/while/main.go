package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {

	i := 0
	for i < 1000 {
		bBool := RandomBool()
		if bBool {
			fmt.Println("檢查成功，執行 ")
		} else {
			fmt.Println("檢查失敗，下一個迴圈, 不會阻塞(wait)！！ ")

		}

		i++
	}
}

func RandomBool() bool {
	return rand.IntN(2) == 0
}
