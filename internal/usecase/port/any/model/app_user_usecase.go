package usecasePortAnyModel

import (
	"example/internal/domain"
)

type AppUserUsecase interface {
	ShowOneByName(sName string) (*domain.AppUser, error)
	ShowOneById(iId uint) (*domain.AppUser, error)
	AddOne(oAppUser *domain.AppUserVariable) error
	IncreaseBalance(iId uint64, iAmount uint64) error
}
