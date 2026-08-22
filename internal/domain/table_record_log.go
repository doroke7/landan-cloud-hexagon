package domain

import (
	"encoding/json"
	"time"
)

type TableRecordLog struct {
	Id            uint            `json:"id"`
	GameId        uint            `json:"game_id"`
	TableRecordId uint            `json:"table_record_id"`
	State         uint8           `json:"state"`
	Text          json.RawMessage `json:"text"`
	Image         json.RawMessage `json:"image"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     time.Time       `json:"deleted_at"`
}

type TableRecordLogRow struct {
	Id            uint            `json:"id"`
	GameId        uint            `json:"game_id"`
	TableRecordId uint            `json:"table_record_id"`
	State         uint8           `json:"state"`
	Text          json.RawMessage `json:"text"`
	Image         json.RawMessage `json:"image"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     time.Time       `json:"deleted_at"`
}

func TableRecordLogRowToTableRecordLog(oRow *TableRecordLogRow) *TableRecordLog {
	return &TableRecordLog{
		Id:            oRow.Id,
		GameId:        oRow.GameId,
		TableRecordId: oRow.TableRecordId,
		State:         oRow.State,
		Text:          oRow.Text,
		Image:         oRow.Image,
		CreatedAt:     oRow.CreatedAt,
		UpdatedAt:     oRow.UpdatedAt,
		DeletedAt:     oRow.DeletedAt,
	}
}

type TableRecordLogValue struct {
	GameId        *uint            `json:"game_id,omitempty"`
	TableRecordId *uint            `json:"table_record_id,omitempty"`
	State         *uint8           `json:"state,omitempty"`
	Text          *json.RawMessage `json:"text,omitempty"`
	Image         *json.RawMessage `json:"image,omitempty"`
	// CreatedAt     *time.Time       `json:"created_at"`
	// UpdatedAt     *time.Time       `json:"updated_at"`
	// DeletedAt     *time.Time       `json:"deleted_at"`
}
