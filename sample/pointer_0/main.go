package main

import "fmt"

func callByPointer1(p *int) {
	*p = 100 // 把傳進來的地址的 那個值改成 100
}

func callByPointer2(p *int) {
	x := 200
	p = &x // 把傳來的地址 改成 x 的
}

func main() {
	a := 1
	b := 2
	callByPointer1(&a)
	callByPointer2(&b) // 這裡發現了一個重點， go 的指標傳入不是傳 ref， 其實只是傳一個地址的值
	// 這裡發現了一個重點， go 的指標傳入不是傳 ref， 其實只是傳一個地址的值
	// 這裡發現了一個重點， go 的指標傳入不是傳 ref， 其實只是傳一個地址的值
	// 這裡發現了一個重點， go 的指標傳入不是傳 ref， 其實只是傳一個地址的值
	// 這裡發現了一個重點， go 的指標傳入不是傳 ref， 其實只是傳一個地址的值
	// 這裡發現了一個重點， go 的指標傳入不是傳 ref， 其實只是傳一個地址的值

	fmt.Println(a)
	fmt.Println(b)

}
