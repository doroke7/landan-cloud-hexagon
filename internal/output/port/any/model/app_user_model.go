package port

import (
	domain "example/internal/domain"
)

type AppUserModel interface {
	IncreaseBalance(iId uint, iAmount uint) (bool, error)
	ShowOneByName(sName string) (*domain.AppUser, error)
	ShowOneById(iId uint) (*domain.AppUser, error)
	AddOne(oAppUser *domain.AppUserValue) (bool, error)
}
