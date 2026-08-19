package main

import "fmt"

// Golang 的 nil pointer 相當 php 的 undefined

type RsaHelper struct {
	Key string
}

func (r *RsaHelper) Encrypt(data string) string {
	return "Encrypted: " + data + " with " + r.Key // 這裡會崩潰，因為 r 是 nil
}

type AbstractMiddleware struct {
	// 注意：這裡是指標 *，預設值是 nil
	rsaHelper *RsaHelper // 如果改成值嵌入就不會報錯
}

func main() {
	// 1. 只宣告，不初始化內部指標
	m := &AbstractMiddleware{}

	// 2. 呼叫方法
	// 在 PHP 中這會報 Fatal Error: Call to a member function Encrypt() on null
	// 在 Go 中這會報 panic: runtime error: invalid memory address or nil pointer dereference
	fmt.Println(m.rsaHelper.Encrypt("secret"))
}
