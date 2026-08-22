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

type TableRow struct {
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

func TableRowToTable(oRow *TableRow) *Table {
	return &Table{
		Id:          oRow.Id,
		No:          oRow.No,
		GameId:      oRow.GameId,
		Key:         oRow.Key,
		State:       oRow.State,
		Description: oRow.Description,
		Result:      oRow.Result,
		StartedAt:   oRow.StartedAt,
		EndedAt:     oRow.EndedAt,
		CreatedAt:   oRow.CreatedAt,
		UpdatedAt:   oRow.UpdatedAt,
		DeletedAt:   oRow.DeletedAt,
		Game:        oRow.Game,
	}
}

// TableDocument 是 mongo 版 TableModel 實際存進去的 schema：Result 在 mysql
// 是 json 欄位，這裡故意存成字串而不是拆成巢狀 bson 文件——Table.Result 是
// json.RawMessage（底層是 []byte），mongo-driver 對 []byte 預設會走 BSON
// binary 編碼，不是我們要的 JSON 文字語意，所以自己控制編解碼，不假手驅動的
// 預設行為。
type TableDocument struct {
	Id          uint      `bson:"_id"`
	No          string    `bson:"no"`
	GameId      uint      `bson:"game_id"`
	Key         string    `bson:"key"`
	State       uint8     `bson:"state"`
	Description string    `bson:"description"`
	Result      string    `bson:"result"`
	StartedAt   time.Time `bson:"started_at"`
	EndedAt     time.Time `bson:"ended_at"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	DeletedAt   time.Time `bson:"deleted_at"`
}

func TableDocumentToTable(oDoc *TableDocument) *Table {
	return &Table{
		Id:          oDoc.Id,
		No:          oDoc.No,
		GameId:      oDoc.GameId,
		Key:         oDoc.Key,
		State:       oDoc.State,
		Description: oDoc.Description,
		Result:      json.RawMessage(oDoc.Result),
		StartedAt:   oDoc.StartedAt,
		EndedAt:     oDoc.EndedAt,
		CreatedAt:   oDoc.CreatedAt,
		UpdatedAt:   oDoc.UpdatedAt,
		DeletedAt:   oDoc.DeletedAt,
	}
}

type TableFilter struct {
	Id     *uint `json:"id,omitempty"`
	GameId *uint `json:"game_id,omitempty"`
}

type TableWhere struct {
	GameId *uint
}
