package main

import "fmt"

type Handler func(int) int

// Middleware 是「包一層再回傳一層」的裝飾器：吃一個 Handler（next），
// 回傳一個新的 Handler。正因為輸入輸出都是 Handler，才能直接巢狀呼叫，
// 例如 Add1Middleware(Multiple2Middleware(Sub2Middleware(fnHandler1)))。
type Middleware func(Handler) Handler

// Chain 把多個 Middleware 疊到 fnHandler 上，aMiddlewares 的順序就是「執行順序」
// （第一個最外層、最先執行），所以由後往前組裝，最後一個先包、第一個最後包。
func Chain(fnHandler Handler, aMiddlewares ...Middleware) Handler {
	for iIndex := len(aMiddlewares) - 1; iIndex >= 0; iIndex-- {
		fnHandler = aMiddlewares[iIndex](fnHandler)
	}
	return fnHandler
}

func Add1Middleware(next Handler) Handler {

	return func(value int) int {

		value = value + 1

		return next(value)
	}
}

func Multiple2Middleware(next Handler) Handler {

	return func(value int) int {

		value = value * 2

		return next(value)
	}
}

func Sub2Middleware(next Handler) Handler {

	return func(value int) int {

		value = value - 2

		return next(value)
	}
}

func Multiple4Middleware(next Handler) Handler {

	return func(value int) int {

		value = value * 4

		return next(value)
	}
}

func main() {

	fnHandler1 := func(value int) int {

		return value
	}

	result := Chain(fnHandler1, Add1Middleware, Multiple2Middleware, Sub2Middleware, Multiple4Middleware)(10)

	fmt.Println(result)

	// 巢狀呼叫組合：每個 XxxMiddleware(next) 回傳的還是一個 Handler，
	// 所以可以直接把上一層的回傳值當下一層的參數，一路巢狀包到最外層 Add1Middleware。
	// 組合完只拿到一個 Handler，還沒真正執行，最後那個 (10) 才是把值帶進去、觸發整條鏈。
	nestedResult2 := Add1Middleware(Multiple2Middleware(Sub2Middleware(fnHandler1)))(10)

	fmt.Println(nestedResult2)

}
