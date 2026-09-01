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

// domain 直接兼作各 adapter 的儲存結構，不再另外開 GameRow（gorm）/
// GameDocument（mongo）兩個欄位完全一樣的鏡像 struct：
//   - gorm：欄位無 tag，靠 NamingStrategy 轉 snake_case；GameType 用 gorm association Preload
//   - mongo：bson tag 指定，_id 對到 Id（counters 累加的 uint）；GameType 用 bson:"-" 略過
//     （Mongo 沒有 join，這個欄位在 mongo 路徑下固定是零值）
type Game struct {
	Id          uint      `json:"id" bson:"_id"`
	GameTypeId  uint      `json:"game_type_id" bson:"game_type_id"`
	Key         string    `json:"key" bson:"key"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" gorm:"default:''" bson:"description"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"default:2038-01-19 03:14:07" bson:"deleted_at"`
	GameType    GameType  `json:"game_type" gorm:"foreignKey:GameTypeId;references:Id" bson:"-"`
}

// TableName 顯式指定表名 games，但 gorm 的 TableName() 是直接取用的原始字串，
// 不會再套用 bootstrap/mysql.go NamingStrategy 設的 TablePrefix，
// 所以這裡自己把 CONFIG.DATABASE.PREFIX 接回去，維持跟 tx-games 一致。
func (Game) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "games"
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

type GameFilter struct {
	Id         *uint `json:"id,omitempty"`
	GameTypeId *uint `json:"game_type_id,omitempty"`
}

type GameWhere struct {
	GameTypeId *uint
}
