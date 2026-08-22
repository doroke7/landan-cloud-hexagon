package domain

import (
	"time"

	bootstrap "example/bootstrap"
)

type GameType struct {
	Id        uint      `json:"id"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// TableName 顯式指定表名 game_types，但 gorm 的 TableName() 是直接取用的原始字串，
// 不會再套用 bootstrap/mysql.go NamingStrategy 設的 TablePrefix，
// 所以這裡自己把 CONFIG.DATABASE.PREFIX 接回去，維持跟 tx-game_types 一致。
func (GameTypeRow) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "game_types"
}

type GameTypeRow struct {
	Id        uint      `json:"id"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func GameTypeRowToGameType(oRow *GameTypeRow) *GameType {
	return &GameType{
		Id:        oRow.Id,
		Key:       oRow.Key,
		Name:      oRow.Name,
		CreatedAt: oRow.CreatedAt,
		UpdatedAt: oRow.UpdatedAt,
		DeletedAt: oRow.DeletedAt,
	}
}

// GameTypeDocument 是 game_type 這個 collection 在 Mongo 裡的儲存結構。
// Mongo 沒有像 SQL 那樣的自增主鍵，_id 改用 counters collection 累加出來的數字，
// 讓 GameType.Id 維持跟其他 adapter 一樣是 uint。
type GameTypeDocument struct {
	Id        uint      `bson:"_id"`
	Key       string    `bson:"key"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

func GameTypeDocumentToGameType(oDoc *GameTypeDocument) *GameType {
	return &GameType{
		Id:        oDoc.Id,
		Key:       oDoc.Key,
		Name:      oDoc.Name,
		CreatedAt: oDoc.CreatedAt,
		UpdatedAt: oDoc.UpdatedAt,
		DeletedAt: oDoc.DeletedAt,
	}
}

type GameTypeValue struct {
	Key  *string `json:"key,omitempty"`
	Name *string `json:"name,omitempty"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
}
