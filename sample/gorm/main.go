package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 對應資料表結構
type User struct {
	ID   uint
	Name string
	Age  int
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 開啟 SQL log（方便觀察）
	db = db.Debug()

	// ===== 查詢範例 =====

	// 1. 查單筆（by ID）
	var user User
	db.First(&user, 1)
	fmt.Println("First:", user)

	// 2. 條件查詢（第一筆）
	db.Where("name = ?", "Alice").First(&user)
	fmt.Println("Where First:", user)

	// 3. 查多筆
	var users []User
	db.Where("age > ?", 20).Find(&users)
	fmt.Println("Find:", users)

	// 4. 只拿一欄（pluck）
	var names []string
	db.Model(&User{}).Where("age > ?", 20).Pluck("name", &names)
	fmt.Println("Names:", names)

	// 5. 計數
	var count int64
	db.Model(&User{}).Where("age > ?", 20).Count(&count)
	fmt.Println("Count:", count)
}
