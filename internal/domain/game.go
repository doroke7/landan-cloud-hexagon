package domain

import (
	"time"
)

type Game struct {
	Id          uint      `json:"id"`
	GameTypeId  uint      `json:"game_type_id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description" gorm:"default:''"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"default:2038-01-19 03:14:07"`
	GameType    GameType  `json:"game_type" gorm:"foreignKey:GameTypeId;references:Id"`
}

type GameRow struct {
	Id          uint      `json:"id"`
	GameTypeId  uint      `json:"game_type_id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description" gorm:"default:''"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"default:2038-01-19 03:14:07"`
	GameType    GameType  `json:"game_type" gorm:"foreignKey:GameTypeId;references:Id"`
}

type GameValue struct {
	GameTypeId  *uint   `json:"game_type_id,omitempty"`
	Key         *string `json:"key,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"` // json:"XXXX,omitempty" 的  omitempty 是告訴 json.Marsal, 如果是 nil 就不解析了
	// CreatedAt   *time.Time `json:"created_at"`
	// UpdatedAt   *time.Time `json:"updated_at"`
	// DeletedAt   *time.Time `json:"deleted_at" gorm:"default:2038-01-19 03:14:07"`
}

func GameRowToGame(oRow *GameRow) *Game {
	return &Game{
		Id:          oRow.Id,
		GameTypeId:  oRow.GameTypeId,
		Key:         oRow.Key,
		Name:        oRow.Name,
		Description: oRow.Description,
		CreatedAt:   oRow.CreatedAt,
		UpdatedAt:   oRow.UpdatedAt,
		DeletedAt:   oRow.DeletedAt,
		GameType:    oRow.GameType,
	}
}

// GameDocument 是 game 這個 collection 在 Mongo 裡的儲存結構。
// Mongo 沒有像 SQL 那樣的自增主鍵，_id 改用 counters collection 累加出來的數字，
// 讓 Game.Id 維持跟其他 adapter 一樣是 uint。
// GameType 沒有做關聯查詢（Mongo 沒有 join），GameDocumentToGame 轉出來的 GameType 固定是零值。
type GameDocument struct {
	Id          uint      `bson:"_id"`
	GameTypeId  uint      `bson:"game_type_id"`
	Key         string    `bson:"key"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	DeletedAt   time.Time `bson:"deleted_at"`
}

func GameDocumentToGame(oDoc *GameDocument) *Game {
	return &Game{
		Id:          oDoc.Id,
		GameTypeId:  oDoc.GameTypeId,
		Key:         oDoc.Key,
		Name:        oDoc.Name,
		Description: oDoc.Description,
		CreatedAt:   oDoc.CreatedAt,
		UpdatedAt:   oDoc.UpdatedAt,
		DeletedAt:   oDoc.DeletedAt,
	}
}

type GameFilter struct {
	Id         *uint `json:"id,omitempty"`
	GameTypeId *uint `json:"game_type_id,omitempty"`
}

type GameWhere struct {
	GameTypeId *uint
}
