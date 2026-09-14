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

type TableRecordLog struct {
	Id            uint64    `json:"id" bson:"_id"`
	GameId        uint64    `json:"game_id" bson:"game_id"`
	TableRecordId uint64    `json:"table_record_id" bson:"table_record_id"`
	State         uint8     `json:"state" bson:"state"`
	Text          string    `json:"text" bson:"text"`
	Image         string    `json:"image" bson:"image"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at" bson:"deleted_at"`
}

type TableRecordLogVariable struct {
	GameId        *uint64 `json:"game_id,omitempty"`
	TableRecordId *uint64 `json:"table_record_id,omitempty"`
	State         *uint8  `json:"state,omitempty"`
	Text          *string `json:"text,omitempty"`
	Image         *string `json:"image,omitempty"`
	// CreatedAt     *time.Time       `json:"created_at"`
	// UpdatedAt     *time.Time       `json:"updated_at"`
	// DeletedAt     *time.Time       `json:"deleted_at"`
}
