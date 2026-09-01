package domain

/*
1. 綜合來說，不建議 proto = domain， 不好改動
2. 建議 保留 gRPC oProtoAdminUser 對 oDomainAdminUser 的轉換，
3. 如果為了節省轉換性能 就 proto domain 合一，增加的性能只有一點，卻增加很大的改動難度
      （譬如， 如果對方的接口格式改變了，我們需要跟著改，那內部domain 也得改）

*/

import (
	"encoding/json"
	"time"

	bootstrap "example/bootstrap"
)

// domain 直接兼作 gorm 的儲存結構，不再另外開欄位完全一樣的 TableRow：
// 欄位無 tag，靠 NamingStrategy 轉 snake_case；Game 用 gorm association Preload。
//
// mongo 那邊還是走獨立的 TableDocument：Result 在 mysql 是 json 欄位，
// json.RawMessage（[]byte）給 mongo-driver 會被編成 BSON binary 而不是 JSON 文字，
// 所以 TableDocument 把它存成 string、自己控制編解碼，沒辦法跟這裡共用同一個 struct。
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

// TableName 顯式指定表名 tables，但 gorm 的 TableName() 是直接取用的原始字串，
// 不會再套用 bootstrap/mysql.go NamingStrategy 設的 TablePrefix，
// 所以這裡自己把 CONFIG.DATABASE.PREFIX 接回去，維持跟 tx-tables 一致。
func (Table) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "tables"
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
