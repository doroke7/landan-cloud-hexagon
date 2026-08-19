package main

import (
	"fmt"
	"reflect"
)

func Add(iX int, iY int) int {
	return iX + iY
}

func Multiply(iX int, iY int) int {
	return iX * iY
}

func Wrapper0(cCaller func(int, int) int) func(int, int) int {
	return cCaller
}

func Wrapper1(cCaller func(int, int) int) func(int, int) int {
	return func(iX int, iY int) int {
		return cCaller(iX, iY)
	}
}

func Wrapper2(cCaller interface{}, args ...interface{}) interface{} {
	// 將傳進來的函數轉為反射對象
	v := reflect.ValueOf(cCaller)

	// 將傳進來的參數 args 轉換為反射需要的 []reflect.Value
	vArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vArgs[i] = reflect.ValueOf(arg)
	}

	// 執行函數 (Call) 並取得回傳值陣列
	results := v.Call(vArgs)

	// 回傳第一個結果
	if len(results) > 0 {
		return results[0].Interface()
	}
	return nil
}

func main() {
	iResult0 := Wrapper0(Add)(1, 2)
	iResult1 := Wrapper1(Add)(1, 2)
	iResult2 := Wrapper2(Add, 1, 2)

	fmt.Printf("Result0: %d\n", iResult0)
	fmt.Printf("Result1: %d\n", iResult1)
	fmt.Printf("Result2: %d\n", iResult2)

}
