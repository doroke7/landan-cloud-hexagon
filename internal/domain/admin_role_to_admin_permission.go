package domain

// AdminRolesToAdminPermission 對應 tx-admin_roles_to_admin_permissions，
// 純樞紐表：只有兩個外鍵、複合主鍵，沒有 id / 時間欄，解綁直接 DELETE。
type AdminRolesToAdminPermission struct {
	AdminRoleId       uint64 `json:"admin_role_id" bson:"admin_role_id" gorm:"primaryKey"`
	AdminPermissionId uint64 `json:"admin_permission_id" bson:"admin_permission_id" gorm:"primaryKey"`
}
