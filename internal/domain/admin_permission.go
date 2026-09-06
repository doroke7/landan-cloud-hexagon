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

type AdminPermission struct {
	Id                     uint64    `json:"id" bson:"_id"`
	AdminPermissionGroupId uint64    `json:"admin_permission_group_id" bson:"admin_permission_group_id"`
	Type                   uint8     `json:"type" bson:"type"` // 1=菜單 2=頁面 3=接口
	Key                    string    `json:"key" bson:"key"`
	Name                   string    `json:"name" bson:"name"`
	CreatedAt              time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt              time.Time `json:"deleted_at" bson:"deleted_at"`
}

type AdminPermissionValue struct {
	Type *uint8  `json:"type,omitempty"`
	Key  *string `json:"key,omitempty"`
	Name *string `json:"name,omitempty"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
}
