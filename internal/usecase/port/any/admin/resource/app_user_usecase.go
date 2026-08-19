package any

import (
// domain "example/internal/domain"
)

type AppUserUsecase interface {
	IncreaseBalance(iId uint, iAmount uint) (bool, error)
}
