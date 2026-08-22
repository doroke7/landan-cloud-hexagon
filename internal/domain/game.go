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

type GameFilter struct {
	Id         *uint `json:"id,omitempty"`
	GameTypeId *uint `json:"game_type_id,omitempty"`
}

type GameWhere struct {
	GameTypeId *uint
}
