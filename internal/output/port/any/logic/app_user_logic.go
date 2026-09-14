package outputPortAnyLogic

import (
	domain "example/internal/domain"
)

type AppUserLogic interface {
	AddAppUser(oAppUser *domain.AppUserVariable) (*domain.AppUser, error)
}
