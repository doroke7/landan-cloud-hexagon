package domain

// AdminUsersToAdminRole 對應 tx-admin_users_to_admin_roles，
// 純樞紐表：只有兩個外鍵、複合主鍵，沒有 id / 時間欄，解綁直接 DELETE。
type AdminUsersToAdminRole struct {
	AdminUserId uint `json:"admin_user_id" bson:"admin_user_id" gorm:"primaryKey"`
	AdminRoleId uint `json:"admin_role_id" bson:"admin_role_id" gorm:"primaryKey"`
}
