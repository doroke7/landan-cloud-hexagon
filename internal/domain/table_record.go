package domain

import (
	"encoding/json"
	"time"

	bootstrap "example/bootstrap"
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

// TableName 顯式指定表名 table_records，但 gorm 的 TableName() 是直接取用的原始字串，
// 不會再套用 bootstrap/mysql.go NamingStrategy 設的 TablePrefix，
// 所以這裡自己把 CONFIG.DATABASE.PREFIX 接回去，維持跟 tx-table_records 一致。
func (TableRecordRow) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "table_records"
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

func TableRecordRowToTableRecord(oRow *TableRecordRow) *TableRecord {
	return &TableRecord{
		Id:        oRow.Id,
		No:        oRow.No,
		GameId:    oRow.GameId,
		TableId:   oRow.TableId,
		State:     oRow.State,
		Text:      oRow.Text,
		Image:     oRow.Image,
		Result:    oRow.Result,
		StartedAt: oRow.StartedAt,
		EndedAt:   oRow.EndedAt,
		CreatedAt: oRow.CreatedAt,
		UpdatedAt: oRow.UpdatedAt,
		DeletedAt: oRow.DeletedAt,
	}
}

// TableRecordDocument 是 mongo 版 TableRecordModel 實際存進去的 schema：Text/Image/Result
// 在 mysql 是 json 欄位，這裡故意存成字串而不是拆成巢狀 bson 文件——它們在 domain 是
// json.RawMessage（底層是 []byte），mongo-driver 對 []byte 預設會走 BSON binary 編碼，
// 不是我們要的 JSON 文字語意，所以自己控制編解碼，不假手驅動的預設行為。
type TableRecordDocument struct {
	Id        uint      `bson:"_id"`
	No        string    `bson:"no"`
	GameId    uint      `bson:"game_id"`
	TableId   uint      `bson:"table_id"`
	State     uint8     `bson:"state"`
	Text      string    `bson:"text"`
	Image     string    `bson:"image"`
	Result    string    `bson:"result"`
	StartedAt time.Time `bson:"started_at"`
	EndedAt   time.Time `bson:"ended_at"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

func TableRecordDocumentToTableRecord(oDoc *TableRecordDocument) *TableRecord {
	return &TableRecord{
		Id:        oDoc.Id,
		No:        oDoc.No,
		GameId:    oDoc.GameId,
		TableId:   oDoc.TableId,
		State:     oDoc.State,
		Text:      json.RawMessage(oDoc.Text),
		Image:     json.RawMessage(oDoc.Image),
		Result:    json.RawMessage(oDoc.Result),
		StartedAt: oDoc.StartedAt,
		EndedAt:   oDoc.EndedAt,
		CreatedAt: oDoc.CreatedAt,
		UpdatedAt: oDoc.UpdatedAt,
		DeletedAt: oDoc.DeletedAt,
	}
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
