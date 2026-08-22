package domain

import bootstrap "example/bootstrap"

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (UserRow) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "users"
}

type UserRow struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserValue struct {
	Name *string `json:"name,omitempty"`
}
