package domain

import bootstrap "example/bootstrap"

type AppUser struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Balance  uint   `json:"balance"`
}

func (AppUserRow) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "app_users"
}

type AppUserRow struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Balance  uint   `json:"balance"`
}

type AppUserValue struct {
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
	Balance  *uint   `json:"balance,omitempty"`
}
