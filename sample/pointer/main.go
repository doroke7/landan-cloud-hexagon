package main

import "fmt"

type Lion struct {
	Name string
}

// 修改函數：接收數值 (Value)
func updateNameByValue(a Lion) {
	a.Name = "Value-Changed" // 這裡改的是副本
}

// 修改函數：接收指針 (Pointer)
func updateNameByPointer(a *Lion) {
	a.Name = "Pointer-Changed" // 這裡改的是原件
}

func main() {

	/*
	  pointer 兩個特性
	  1. 函數操作輸入時候：
	       輸入值的時候-> 修改副本；
	       輸入指針-> 修改原來實例。

	  2. 變量賦予的時候，修改變量：
	       值的會是拷貝 -》只會改當下變量
	       指針的會是相同內存地址 -》 會統一修改

	*/

	oLion1 := Lion{Name: "Original-1"}
	oLion22 := Lion{Name: "Original-11"}

	pLion := &Lion{Name: "Original-2"}

	fmt.Println("--- 修改前 ---")
	fmt.Printf("oLion: %v\n", oLion1.Name)
	fmt.Printf("pLion: %v\n", pLion.Name)
	fmt.Printf("oLion22: %v\n", oLion22.Name)

	// 嘗試修改
	updateNameByValue(oLion1)  // 傳入的是 oLion 的整塊內存拷貝
	updateNameByPointer(pLion) // 傳入的是 pLion 的內存地址

	fmt.Println("\n--- 修改後 ---")
	fmt.Printf("oLion: %v (沒變，因為函數內改的是副本)\n", oLion1.Name)
	fmt.Printf("pLion: %v (變了，因為透過地址直接改了原件)\n", pLion.Name)

	pLion2 := pLion // o3 現在和 pLion 指向同一個地址
	pLion2.Name = "Vibing"
	fmt.Printf("\npLion 被 o3 修改後: %v\n", pLion.Name)

}
