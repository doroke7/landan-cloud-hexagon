package main

import "fmt"

// Response 是個泛型結構體，[T any] 定義了它的「型別參數」
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"` // 這裡的 T 會根據對接的數據改變
}

type User struct {
	ID   uint64
	Name string
}

func main() {
	// 1 (人的努力)：對接到 User 型別
	userRes := Response[User]{
		Code: 200,
		Data: User{ID: 1, Name: "Senior Architect"},
	}

	// 1 (人的努力)：對接到 string 型別
	msgRes := Response[string]{
		Code: 200,
		Data: "Operation Successful",
	}

	fmt.Printf("User: %+v\n", userRes.Data.Name)
	fmt.Printf("Message: %s\n", msgRes.Data)
}
