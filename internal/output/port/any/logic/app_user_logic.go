package logic

import (
	domain "example/internal/domain"
)

type AppUserLogic interface {
	AddAppUser(oAppUser *domain.AppUserValue) (*domain.AppUser, error)
}
