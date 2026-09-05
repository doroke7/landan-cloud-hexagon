package usecasePortAnyAdminResource

import (
// domain "example/internal/domain"
)

type AppUserUsecase interface {
	IncreaseBalance(iId uint64, iAmount uint64) error
}
