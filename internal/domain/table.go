package domain

/*
1. 綜合來說，不建議 proto = domain， 不好改動
2. 建議 保留 gRPC oProtoAdminUser 對 oDomainAdminUser 的轉換，
3. 如果為了節省轉換性能 就 proto domain 合一，增加的性能只有一點，卻增加很大的改動難度
      （譬如， 如果對方的接口格式改變了，我們需要跟著改，那內部domain 也得改）

*/

import (
	"time"

	bootstrap "example/bootstrap"
)

// domain 直接兼作各 adapter 的儲存結構，不再另外開 TableRow（gorm）/
// TableDocument（mongo）兩個欄位完全一樣的鏡像 struct：
//   - gorm：欄位無 tag，靠 NamingStrategy 轉 snake_case；Game 用 gorm association Preload
//   - mongo：bson tag 指定，_id 對到 Id（counters 累加的 uint）；Game 用 bson:"-" 略過
//     （Mongo 沒有 join，這個欄位在 mongo 路徑下固定是零值）
//
// Result 是純字串（不是 json.RawMessage），mongo 直接存字串、不拆成巢狀 bson 文件。
type Table struct {
	Id          uint      `json:"id" bson:"_id"`
	No          string    `json:"no" bson:"no"`
	GameId      uint      `json:"game_id" bson:"game_id"`
	Key         string    `json:"key" bson:"key"`
	State       uint8     `json:"state" bson:"state"`
	Description string    `json:"description" bson:"description"`
	Result      string    `json:"result" bson:"result"`
	StartedAt   time.Time `json:"started_at" bson:"started_at"`
	EndedAt     time.Time `json:"ended_at" bson:"ended_at"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at" bson:"deleted_at"`
	Game        Game      `json:"game" gorm:"foreignKey:GameId;references:Id" bson:"-"`
}

// TableName 顯式指定表名 tables，但 gorm 的 TableName() 是直接取用的原始字串，
// 不會再套用 bootstrap/mysql.go NamingStrategy 設的 TablePrefix，
// 所以這裡自己把 CONFIG.DATABASE.PREFIX 接回去，維持跟 tx-tables 一致。
func (Table) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "tables"
}

type TableValue struct {
	No          *string    `json:"no,omitempty"`
	GameId      *uint      `json:"game_id,omitempty"`
	Key         *string    `json:"key,omitempty"`
	State       *uint8     `json:"state,omitempty"`
	Description *string    `json:"description,omitempty"`
	Result      *string    `json:"result,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`
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
