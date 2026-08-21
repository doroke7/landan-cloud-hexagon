package domain

import (
	"encoding/json"
	"time"
)

type TableRecord struct {
	Id        uint            `json:"id"`
	No        string          `json:"no"` // 年-月日-桌號-局號-期號
	GameId    uint            `json:"game_id"`
	TableId   uint            `json:"table_id"`
	State     uint8           `json:"state"`
	Text      json.RawMessage `json:"text"`
	Image     json.RawMessage `json:"image"`
	Result    json.RawMessage `json:"result"`
	StartedAt time.Time       `json:"started_at"`
	EndedAt   time.Time       `json:"ended_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt time.Time       `json:"deleted_at"`
}

type TableRecordRow struct {
	Id        uint            `json:"id"`
	No        string          `json:"no"` // 年-月日-桌號-局號-期號
	GameId    uint            `json:"game_id"`
	TableId   uint            `json:"table_id"`
	State     uint8           `json:"state"`
	Text      json.RawMessage `json:"text"`
	Image     json.RawMessage `json:"image"`
	Result    json.RawMessage `json:"result"`
	StartedAt time.Time       `json:"started_at"`
	EndedAt   time.Time       `json:"ended_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt time.Time       `json:"deleted_at"`
}

type TableRecordValue struct {
	No        *string          `json:"no,omitempty"` // 年-月日-桌號-局號-期號
	GameId    *uint            `json:"game_id,omitempty"`
	TableId   *uint            `json:"table_id,omitempty"`
	State     *uint8           `json:"state,omitempty"`
	Text      *json.RawMessage `json:"text,omitempty"`
	Image     *json.RawMessage `json:"image,omitempty"`
	Result    *json.RawMessage `json:"result,omitempty"`
	StartedAt *time.Time       `json:"started_at,omitempty"`
	EndedAt   *time.Time       `json:"ended_at,omitempty"`
	// CreatedAt *time.Time       `json:"created_at"`
	// UpdatedAt *time.Time       `json:"updated_at"`
	// DeletedAt *time.Time       `json:"deleted_at"`
}
