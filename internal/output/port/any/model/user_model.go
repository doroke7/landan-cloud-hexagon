package outputPortAnyModel

import (
	domain "example/internal/domain"
)

type UserModel interface {
	AddOne(oUser *domain.UserVariable) error
	ShowOneById(iId uint64) (*domain.User, error)
}
