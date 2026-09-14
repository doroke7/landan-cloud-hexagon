package domain

type AdminRolesToAdminPermission struct {
	AdminRoleId       uint64 `json:"admin_role_id" bson:"admin_role_id" gorm:"primaryKey"`
	AdminPermissionId uint64 `json:"admin_permission_id" bson:"admin_permission_id" gorm:"primaryKey"`
}
