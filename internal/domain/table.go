package domain

import (
	"encoding/json"
	"time"
)

type Table struct {
	Id          uint            `json:"id"`
	No          string          `json:"no"`
	GameId      uint            `json:"game_id"`
	Key         string          `json:"key"`
	State       uint8           `json:"state"`
	Description string          `json:"description"`
	Result      json.RawMessage `json:"result"`
	StartedAt   time.Time       `json:"started_at"`
	EndedAt     time.Time       `json:"ended_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   time.Time       `json:"deleted_at"`
	Game        Game            `json:"game" gorm:"foreignKey:GameId;references:Id"`
}

type TableValue struct {
	No          *string          `json:"no,omitempty"`
	GameId      *uint            `json:"game_id,omitempty"`
	Key         *string          `json:"key,omitempty"`
	State       *uint8           `json:"state,omitempty"`
	Description *string          `json:"description,omitempty"`
	Result      *json.RawMessage `json:"result,omitempty"`
	StartedAt   *time.Time       `json:"started_at,omitempty"`
	EndedAt     *time.Time       `json:"ended_at,omitempty"`
	// CreatedAt   *time.Time       `json:"created_at"`
	// UpdatedAt   *time.Time       `json:"updated_at"`
	// DeletedAt   *time.Time       `json:"deleted_at"`
}

type TableFilter struct {
	Id     *uint `json:"id,omitempty"`
	GameId *uint `json:"game_id,omitempty"`
}

type TableWhere struct {
	GameId *uint
}
