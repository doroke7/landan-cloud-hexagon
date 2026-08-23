package usecasePortAnyModel

import (
	"example/internal/domain"
)

type AppUserUsecase interface {
	ShowOneByName(sName string) (*domain.AppUser, error)
	ShowOneById(iId uint) (*domain.AppUser, error)
	AddOne(oAppUser *domain.AppUserValue) (bool, error)
	IncreaseBalance(iId uint, iAmount uint) (bool, error)
}
