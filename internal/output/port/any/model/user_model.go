package outputPortAnyModel

import (
	domain "example/internal/domain"
)

type UserModel interface {
	AddOne(oUser *domain.UserValue) error
	ShowOneById(iId int) (*domain.User, error)
}
