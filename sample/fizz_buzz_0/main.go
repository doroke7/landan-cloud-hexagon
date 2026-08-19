package main

import "fmt"

/*


Fizz Buzz 是一個在軟體工程師面試中極其著名的經典基礎問題，經常用來作為撰寫程式的初階篩選（Screening test）。
它的描述非常簡單，但如果放到**併發編程（Concurrency）**的時空背景下，它會演變成一個非常考驗同步技巧的進階難題（也就是著名的 多執行緒 Fizz Buzz 問題）。
我們先來看最基礎的版本，再看它如何變成併發痛點。


*/

func FizzBuzz(iNumber int) {
	for i := 1; i <= iNumber; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Println("FizzBuzz")
		} else if i%3 == 0 {
			fmt.Println("Fizz")
		} else if i%5 == 0 {
			fmt.Println("Buzz")
		} else {
			fmt.Println(i)
		}
	}
}

func main() {

	FizzBuzz(100)

}
