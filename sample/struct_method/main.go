package main

import "fmt"

type User struct {
	Name string
	Age  int
}

// 【值綁定】
// 操作 值綁定 UpdateNameValue 的時候，
// 函數裡面的 oUser 會從新拷貝一個新的 User 值， 裡面外面的 user 不會干擾

func (oUser User) UpdateNameValue(newName string) {
	oUser.Name = newName
	fmt.Printf("  >> 函數內部（值綁定）：%s\n", oUser.Name)
	fmt.Printf("  >> 函數內部（值綁定）：%v\n", oUser)

}

// 【指標綁定】
// 操作 值綁定 UpdateNamePointer 的時候，
// 函數裡面的 oUser 會從新拷貝地址，外面裡面的 User 相同

func (oUser *User) UpdateNamePointer(newName string) {
	oUser.Name = newName
	fmt.Printf("  >> 函數內部（指標綁定）：%s\n", oUser.Name)
}

func main() {
	// 初始化一個 Landan
	oUser1 := User{
		Name: "Landan",
		Age:  30,
	}

	fmt.Println("1. 嘗試使用【值綁定】修改：")
	oUser1.UpdateNameValue("XYZ")
	fmt.Printf("   結果：原本的 user 名字還是 %s (沒變！)\n\n", oUser1.Name)

	// fmt.Println("2. 嘗試使用【指標綁定】修改：")
	// oUser1.UpdateNamePointer("ABC")
	// fmt.Printf("   結果：原本的 user 名字變成了 %s (成功！)\n", oUser1.Name)
}
