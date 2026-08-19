package main

import "fmt"

type Handler func(int) int

// Middleware 現在直接吃 value 跟 next，不再是「func(Handler) Handler」包一層再回傳一層，
// next 就是內側函數的參數，跟 value 平行放，middleware 本身不用自己再包 closure。
type Middleware func(value int, next Handler) int

// Chain 把多個 Middleware 疊到 fnHandler 上，aMiddlewares 的順序就是「執行順序」
// （第一個最外層、最先執行）。組裝 closure 的動作集中在這裡做一次，
// 每個 Middleware 自己只是單純的 func(value, next) int，不用各自處理 closure。
func Chain(fnHandler Handler, aMiddlewares ...Middleware) Handler {
	for iIndex := len(aMiddlewares) - 1; iIndex >= 0; iIndex-- {
		fnMiddleware := aMiddlewares[iIndex]
		fnNext := fnHandler
		fnHandler = func(value int) int {
			return fnMiddleware(value, fnNext)
		}
	}
	return fnHandler
}

func Add1Middleware(value int, next Handler) int {

	value = value + 1

	return next(value)
}

func Multiple2Middleware(value int, next Handler) int {

	value = value * 2

	return next(value)
}

func Sub2Middleware(value int, next Handler) int {

	value = value - 2

	return next(value)
}

func Multiple4Middleware(value int, next Handler) int {

	value = value * 4

	return next(value)
}

func main() {

	fnHandler1 := func(value int) int {

		return value
	}

	fnHandler5 := Chain(fnHandler1, Add1Middleware, Multiple2Middleware, Sub2Middleware, Multiple4Middleware)

	result := fnHandler5(10)

	fmt.Println(result)

	// 不用 Chain，手動一層層包：Chain 做的事情就是這個 for 迴圈，
	// 這裡展開來寫，效果跟上面 Chain(...) 完全一樣。
	fnManualHandler4 := func(value int) int {
		return Multiple4Middleware(value, fnHandler1)
	}
	fnManualHandler3 := func(value int) int {
		return Sub2Middleware(value, fnManualHandler4)
	}
	fnManualHandler2 := func(value int) int {
		return Multiple2Middleware(value, fnManualHandler3)
	}
	fnManualHandler1 := func(value int) int {
		return Add1Middleware(value, fnManualHandler2)
	}

	manualResult := fnManualHandler1(10)

	fmt.Println(manualResult)

	nestedResult := Add1Middleware(10, func(value int) int {
		return Multiple2Middleware(value, func(value int) int {
			return Sub2Middleware(value, func(value int) int {
				return Multiple4Middleware(value, fnHandler1)
			})
		})
	})

	fmt.Println(nestedResult)

	nestedResult2 := Add1Middleware(10, func(value int) int {
		return Multiple2Middleware(value, fnHandler1)
	})

	fmt.Println(nestedResult2)

}
