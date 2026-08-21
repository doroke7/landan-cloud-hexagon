package domain

import (
	"time"
)

type AdminUser struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type AdminUserRow struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type AdminUserValue struct {
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
	// CreatedAt *time.Time `json:"created_at"`
	// UpdatedAt *time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at"`
}

func AdminUserRowToAdminUser(oRow *AdminUserRow) *AdminUser {
	return &AdminUser{
		Id:        oRow.Id,
		Name:      oRow.Name,
		Password:  oRow.Password,
		CreatedAt: oRow.CreatedAt,
		UpdatedAt: oRow.UpdatedAt,
		DeletedAt: oRow.DeletedAt,
	}
}
