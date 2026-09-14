package outputPortAnyModel

import (
	domain "example/internal/domain"
)

type AppUserModel interface {
	IncreaseBalance(iId uint, iAmount uint64) error
	ShowOneByName(sName string) (*domain.AppUser, error)
	ShowOneById(iId uint) (*domain.AppUser, error)
	AddOne(oAppUser *domain.AppUserVariable) error
}
