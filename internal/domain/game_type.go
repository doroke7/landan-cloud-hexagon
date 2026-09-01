package domain

/*
1. 綜合來說，不建議 proto = domain， 不好改動
2. 建議 保留 gRPC oProtoAdminUser 對 oDomainAdminUser 的轉換，
3. 如果為了節省轉換性能 就 proto domain 合一，增加的性能只有一點，卻增加很大的改動難度
      （譬如， 如果對方的接口格式改變了，我們需要跟著改，那內部domain 也得改）

*/

import (
	"time"
)

// domain 直接兼作各 adapter 的儲存結構，不再另外開 GameTypeRow（gorm）/
// GameTypeDocument（mongo）兩個欄位完全一樣的鏡像 struct：
//   - gorm：欄位無 tag，靠 NamingStrategy 轉 snake_case；主鍵 Id -> id
//   - mongo：bson tag 指定，_id 對到 Id（counters 累加的 uint，不是 ObjectID）
type GameType struct {
	Id        uint      `json:"id" bson:"_id"`
	Key       string    `json:"key" bson:"key"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at"`
}

type GameTypeValue struct {
	Key  *string `json:"key,omitempty"`
	Name *string `json:"name,omitempty"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
}
