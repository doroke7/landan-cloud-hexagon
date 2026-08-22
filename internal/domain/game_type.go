package domain

import (
	"time"
)

type GameType struct {
	Id        uint      `json:"id"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
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

type GameTypeValue struct {
	Key  *string `json:"key,omitempty"`
	Name *string `json:"name,omitempty"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
}
