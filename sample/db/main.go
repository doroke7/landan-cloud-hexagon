package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func main() {
	Db()
}

func Db() *gorm.DB {
	// 1. 直接寫死你的 Docker MySQL 連線字串 (DSN)
	// 格式: user:password@tcp(container_name:3306)/dbname?params
	masterDSN := "backend:mysql_pass_backend@tcp(mysql:3306)/gogin?charset=utf8mb4&parseTime=True&loc=Local"
	slaveDSN := "backend:mysql_pass_backend@tcp(mysql:3306)/gogin?charset=utf8mb4&parseTime=True&loc=Local"

	fmt.Println("正在連線至:", masterDSN)

	// 2. 初始化 GORM
	// 注意：這裡加入了錯誤檢查，不要用 "_"，否則連不上你也看不出來
	oDb, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "gin-", // 直接寫死前綴
		},
	})

	if err != nil {
		// 如果這裡噴掉，代表你的 Go 容器真的連不到 MySQL 容器
		panic(fmt.Sprintf("資料庫連線失敗: %v", err))
	}

	// 3. 註冊讀寫分離 (雖然目前都連同一台，但結構要留著)
	err = oDb.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(masterDSN)},
		Replicas: []gorm.Dialector{mysql.Open(slaveDSN)},
		Policy:   dbresolver.RandomPolicy{},
	}).
		SetMaxIdleConns(10).
		SetConnMaxLifetime(time.Hour))

	if err != nil {
		panic(fmt.Sprintf("DBResolver 註冊失敗: %v", err))
	}

	// 驗證連線是否真的活著
	sqlDB, _ := oDb.DB()
	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("資料庫 Ping 失敗 (雖然 Open 成功了): %v", err))
	}

	fmt.Println("oDb 初始化成功:", oDb)
	return oDb
}
