package domain

import (
	"encoding/json"
	"time"

	bootstrap "example/bootstrap"
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

func (TableRecordLogRow) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "table_record_logs"
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

// TableRecordLogDocument 是 table_record_log 這個 collection 在 Mongo 裡的儲存結構。
// Mongo 沒有像 SQL 那樣的自增主鍵，_id 改用 counters collection 累加出來的數字，
// 讓 TableRecordLog.Id 維持跟其他 adapter 一樣是 uint。
// Text/Image 故意存成字串而不是拆成巢狀 bson 文件，理由跟 TableDocument.Result 一樣。
type TableRecordLogDocument struct {
	Id            uint      `bson:"_id"`
	GameId        uint      `bson:"game_id"`
	TableRecordId uint      `bson:"table_record_id"`
	State         uint8     `bson:"state"`
	Text          string    `bson:"text"`
	Image         string    `bson:"image"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
	DeletedAt     time.Time `bson:"deleted_at"`
}

func TableRecordLogDocumentToTableRecordLog(oDoc *TableRecordLogDocument) *TableRecordLog {
	return &TableRecordLog{
		Id:            oDoc.Id,
		GameId:        oDoc.GameId,
		TableRecordId: oDoc.TableRecordId,
		State:         oDoc.State,
		Text:          json.RawMessage(oDoc.Text),
		Image:         json.RawMessage(oDoc.Image),
		CreatedAt:     oDoc.CreatedAt,
		UpdatedAt:     oDoc.UpdatedAt,
		DeletedAt:     oDoc.DeletedAt,
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
