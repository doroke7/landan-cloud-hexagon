package any

import (
	"example/internal/domain"
)

type AppUserUsecase interface {
	AddAppUser(oAppUser *domain.AppUser) (bool, error)
}
