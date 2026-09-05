package usecasePortAnyLogic

import (
	"example/internal/domain"
)

type AppUserUsecase interface {
	AddAppUser(oAppUser *domain.AppUser) error
}
