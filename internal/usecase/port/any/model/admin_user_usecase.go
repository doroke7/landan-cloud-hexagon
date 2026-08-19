package any

import (
	"example/internal/domain"
)

type AdminUserUsecase interface {
	ShowOneByName(sName string) (*domain.AdminUser, error)
	ShowOneById(iId uint) (*domain.AdminUser, error)
	AddOne(oAdminUser *domain.AdminUserValue) (bool, error)
}
