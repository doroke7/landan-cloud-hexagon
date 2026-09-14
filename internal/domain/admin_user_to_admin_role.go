package domain

type AdminUsersToAdminRole struct {
	AdminUserId uint64 `json:"admin_user_id" bson:"admin_user_id" gorm:"primaryKey"`
	AdminRoleId uint64 `json:"admin_role_id" bson:"admin_role_id" gorm:"primaryKey"`
}
