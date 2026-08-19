package main

import "fmt"

/*
公式： 函數名稱開槽T的任意型別，函數輸入使用了這個變動型別， 輸出的時候可能是這個型別

func 函數名稱[T any](aRows []T, func(T) bool) []T {

}


*/

func Filter[T any](slice []T, test func(T) bool) []T {
	var result []T
	for _, item := range slice {
		// 執行外部注入的判斷條件
		if test(item) {
			result = append(result, item)
		}
	}
	return result
}

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	// 數據源
	users := []User{
		{ID: 1, Name: "張三", Age: 18},
		{ID: 2, Name: "李四", Age: 30},
		{ID: 3, Name: "王五", Age: 25},
	}

	// 在函數使用這個泛型：對接到「年紀大於 20」的邏輯
	// 注意：調用時不需要寫 Filter[User]，Go 會自動推導 (Type Inference)
	adults := Filter(users, func(u User) bool {
		return u.Age > 20
	})

	fmt.Printf("成年用戶數: %d, 第一位是: %s\n", len(adults), adults[0].Name)

	// 同一個函數，對接到「純數字」過濾
	nums := []int{1, 2, 3, 4, 5, 6}
	evenNums := Filter(nums, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println("偶數列表:", evenNums)
}
