package main

type User struct {
	Id   uint
	Name string

	Detail UserDetail `gorm:"foreignKey:UserId;references:Id"`
}

type UserDetail struct {
	Id uint

	UserId uint
	Age    int
}

/*
RDBMS 的 foreignKey 指的是，兩個表關係中，
使用關係聯動的 「存放另一張表 Key 值的那個欄位」，就是外鍵
跟 關係的方向 無關。

所以不管從 A表，還是B表去看， 外鍵都是同一個

*/
//////////////////////////////////////////////////////////////////////////////

type Game struct {
	Id         uint
	Name       string
	GameTypeId uint

	GameType GameType `gorm:"foreignKey:GameTypeId;references:Id"`
}

type GameType struct {
	Id   uint
	Name string
}
