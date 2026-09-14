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
)

type AdminPermissionGroup struct {
	Id               uint64                  `json:"id" bson:"_id"`
	ParentId         uint64                  `json:"parent_id" bson:"parent_id"`
	Key              string                  `json:"key" bson:"key"`
	Name             string                  `json:"name" bson:"name"`
	Parent           *AdminPermissionGroup   `json:"parent,omitempty" bson:"-" gorm:"foreignKey:ParentId;references:Id"`
	Children         []*AdminPermissionGroup `json:"children,omitempty" bson:"-" gorm:"foreignKey:ParentId"`
	AdminPermissions []*AdminPermission      `json:"admin_permissions,omitempty" bson:"-" gorm:"foreignKey:AdminPermissionGroupId"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at"`
}

type AdminPermissionGroupVariable struct {
	Key              *string                    `json:"key,omitempty"`
	Name             *string                    `json:"name,omitempty"`
	AdminPermissions []*AdminPermissionVariable `json:"admin_permissions" bson:"-" `
	ParentId         *uint64                    `json:"parent_id,omitempty"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
}

func (oSelf AdminPermissionGroup) MarshalJSON() ([]byte, error) {
	type alias AdminPermissionGroup

	oCopy := alias(oSelf)
	if oCopy.Children == nil {
		oCopy.Children = []*AdminPermissionGroup{}
	}
	if oCopy.AdminPermissions == nil {
		oCopy.AdminPermissions = []*AdminPermission{}
	}

	aBytes, oErr := json.Marshal(oCopy)

	return aBytes, oErr
}
