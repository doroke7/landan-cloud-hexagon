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

type AdminUser struct {
	Id         uint64      `json:"id" bson:"_id"`
	Name       string      `json:"name" bson:"name"`
	Password   string      `json:"password" bson:"password"`
	CreatedAt  time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" bson:"updated_at"`
	DeletedAt  time.Time   `json:"deleted_at" bson:"deleted_at"`
	AdminRoles []AdminRole `json:"admin_roles,omitempty" bson:"-" gorm:"many2many:admin_users_to_admin_roles"`
}

type AdminUserVariable struct {
	Name         *string  `json:"name,omitempty" validate:"omitempty,min=4"`
	Password     *string  `json:"password,omitempty" validate:"omitempty,min=8"`
	AdminRoleIds []uint64 `json:"admin_role_ids,omitempty" validate:"omitempty,dive,gt=0"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
}
