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

type GameType struct {
	Id        uint       `json:"id" bson:"_id"`
	ParentId  uint       `json:"parent_id" bson:"parent_id"`
	Key       string     `json:"key" bson:"key"`
	Name      string     `json:"name" bson:"name"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time  `json:"deleted_at" bson:"deleted_at"`
	Children  []GameType `gorm:"foreignKey:ParentId"`
}

type GameTypeValue struct {
	Key      *string `json:"key,omitempty"`
	ParentId *uint   `json:"parent_id,omitempty"`
	Name     *string `json:"name,omitempty"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
}
