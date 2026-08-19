package main

import (
	"fmt"
	"strconv"
	"sync"
)

type AppUserModel struct {
	Id       int
	Nickname string
	// 將 Once 放在 struct 內部
	once sync.Once
}

/*
將 sync.Once 放進 struct 的方法中，要保證「該實例的某個 method 只執行一次」，
核心技巧在於：將 sync.Once 對象直接作為 struct 的一個成員欄位。
這樣，這把「鎖」就會跟著實例走。無論這個實例的方法被併發調用了多少次，只要是同一個對象，它就只會執行一次。
*/

// 2. 不要讓 Once 逃逸
// sync.Once 是不能被複製的（它內部含有鎖）。
// 一旦 struct 被初始化並開始使用，就不要再對這個 user 對象進行賦值拷貝（例如 user2 := *user1），
// 否則會引發意想不到的併發錯誤。

func (oSelf *AppUserModel) Do(sString string) {
	// Self.once 是跟著實例走的，也就是說 同一個實例確保只會執行一次。
	oSelf.once.Do(func() {
		oSelf.Nickname = "Landan_Makati:" + sString
		fmt.Printf("===> [ID:%d] 正在執行昂貴的初始化（僅限一次）AppUserModel.Nickname= %s\n", oSelf.Id, oSelf.Nickname)

	})
}

func main() {
	// 建立一個實例
	oAppUser1 := &AppUserModel{
		Id: 1,
	}
	oAppUser2 := &AppUserModel{
		Id: 2,
	}
	var oWaitGroup1 sync.WaitGroup
	var oWaitGroup2 sync.WaitGroup

	// 模擬多個併發請求同時調用同一個 user 的方法
	for i := 0; i < 5000; i++ {
		oWaitGroup1.Add(1)
		go func(n int) {
			defer oWaitGroup1.Done()
			sI := strconv.Itoa(i)
			// fmt.Printf("Goroutine %d 嘗試調用...\n", n)
			oAppUser1.Do(sI)
		}(i)
	}

	for i := 0; i < 5000; i++ {
		oWaitGroup2.Add(1)
		go func(n int) {
			defer oWaitGroup2.Done()
			sI := strconv.Itoa(i)
			// fmt.Printf("Goroutine %d 嘗試調用...\n", n)
			oAppUser2.Do(sI)
		}(i)
	}

	oWaitGroup1.Wait()
	oWaitGroup2.Wait()

	fmt.Printf("最終結果: oAppUser1 %s\n", oAppUser1.Nickname)
	fmt.Printf("最終結果: oAppUser2 %s\n", oAppUser2.Nickname)

}
