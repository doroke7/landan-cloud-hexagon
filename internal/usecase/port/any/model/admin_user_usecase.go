package usecasePortAnyModel

import (
	"example/internal/domain"
)

type AdminUserUsecase interface {
	ShowOneByName(sName string) (*domain.AdminUser, error)
	ShowOneById(iId uint) (*domain.AdminUser, error)
	AddOne(oAdminUser *domain.AdminUserValue) (bool, error)
	EditOneById(oAdminUser *domain.AdminUserValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)
}
