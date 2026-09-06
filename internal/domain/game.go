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

type Game struct {
	Id          uint64    `json:"id" bson:"_id"`
	GameTypeId  uint64    `json:"game_type_id" bson:"game_type_id"`
	Key         string    `json:"key" bson:"key"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" gorm:"default:''" bson:"description"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"default:2038-01-19 03:14:07" bson:"deleted_at"`
	GameType    GameType  `json:"game_type" gorm:"foreignKey:GameTypeId;references:Id" bson:"-"`
}

type GameValue struct {
	GameTypeId  *uint64 `json:"game_type_id,omitempty"`
	Key         *string `json:"key,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"` // json:"XXXX,omitempty" 的  omitempty 是告訴 json.Marsal, 如果是 nil 就不解析了
	// CreatedAt   *time.Time `json:"created_at"`
	// UpdatedAt   *time.Time `json:"updated_at"`
	// DeletedAt   *time.Time `json:"deleted_at" gorm:"default:2038-01-19 03:14:07"`
}
