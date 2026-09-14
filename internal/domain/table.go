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

type Table struct {
	Id          uint64    `json:"id" bson:"_id"`
	No          string    `json:"no" bson:"no"`
	GameId      uint64    `json:"game_id" bson:"game_id"`
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

type TableVariable struct {
	No          *string    `json:"no,omitempty"`
	GameId      *uint64    `json:"game_id,omitempty"`
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
