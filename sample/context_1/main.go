package main

import (
	"context"
	"fmt"
)

// authorizationKey 是 context key 的型別，用空 struct 而不是字串，避免跟別人塞的
// key 撞名——這是 Go 官方建議的 context key 慣例。
type authorizationKey struct{}

// WithAuthorization 建一個空的 *string holder 塞進 context，回傳新的 context。
// 要在呼叫下游之前做，下游才能透過 Set 寫進同一個 holder。
func WithAuthorization(oContext context.Context) context.Context {
	return context.WithValue(oContext, authorizationKey{}, new(string))
}

// Set 把值寫進 context 裡的 holder；如果 context 沒有先呼叫過 WithAuthorization
// （holder 不存在），回傳 false，不會 panic。
func Set(oContext context.Context, sValue string) bool {
	pValue, bOk := oContext.Value(authorizationKey{}).(*string)
	if !bOk {
		return false
	}
	*pValue = sValue
	return true
}

// Get 讀出 context 裡 holder 目前的值。
func Get(oContext context.Context) (string, bool) {
	pValue, bOk := oContext.Value(authorizationKey{}).(*string)
	if !bOk {
		return "", false
	}
	return *pValue, true
}

func main() {
	oContext, fCancel := context.WithCancel(context.Background())
	defer fCancel()

	oContext = WithAuthorization(oContext)

	test(oContext)

	sAuthorization, bOk := Get(oContext)
	fmt.Println("sAuthorization =", sAuthorization, "ok =", bOk)
}

func test(oContext context.Context) {
	Set(oContext, "12345678")
}
