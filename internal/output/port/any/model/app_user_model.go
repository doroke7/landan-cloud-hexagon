package port

import (
// domain "example/internal/domain"
)

type AppUserModel interface {
	IncreaseBalance(iId uint, iAmount uint) (bool, error)
}
