package domain

/*
1. 綜合來說，不建議 proto = domain， 不好改動
2. 建議 保留 gRPC oProtoAdminUser 對 oDomainAdminUser 的轉換，
3. 如果為了節省轉換性能 就 proto domain 合一，增加的性能只有一點，卻增加很大的改動難度
      （譬如， 如果對方的接口格式改變了，我們需要跟著改，那內部domain 也得改）

*/

import (
	"time"

	bootstrap "example/bootstrap"
)

type AdminUser struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// TableName 顯式指定表名 admin_users，但 gorm 的 TableName() 是直接取用的原始字串，
// 不會再套用 bootstrap/mysql.go NamingStrategy 設的 TablePrefix，
// 所以這裡自己把 CONFIG.DATABASE.PREFIX 接回去，維持跟 tx-admin_users 一致。
func (AdminUserRow) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "admin_users"
}

type AdminUserRow struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type AdminUserValue struct {
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
}

func AdminUserRowToAdminUser(oRow *AdminUserRow) *AdminUser {
	return &AdminUser{
		Id:        oRow.Id,
		Name:      oRow.Name,
		Password:  oRow.Password,
		CreatedAt: oRow.CreatedAt,
		UpdatedAt: oRow.UpdatedAt,
		DeletedAt: oRow.DeletedAt,
	}
}

// AdminUserDocument 是 admin_user 這個 collection 在 Mongo 裡的儲存結構。
// Mongo 沒有像 SQL 那樣的自增主鍵，_id 改用 counters collection 累加出來的數字，
// 讓 AdminUser.Id 維持跟其他 adapter 一樣是 uint。
type AdminUserDocument struct {
	Id        uint      `bson:"_id"`
	Name      string    `bson:"name"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

func AdminUserDocumentToAdminUser(oDoc *AdminUserDocument) *AdminUser {
	return &AdminUser{
		Id:        oDoc.Id,
		Name:      oDoc.Name,
		Password:  oDoc.Password,
		CreatedAt: oDoc.CreatedAt,
		UpdatedAt: oDoc.UpdatedAt,
		DeletedAt: oDoc.DeletedAt,
	}
}
