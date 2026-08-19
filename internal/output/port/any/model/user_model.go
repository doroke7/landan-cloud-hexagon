package port

import (
	domain "example/internal/domain"
)

type UserModel interface {
	AddOne(oUser *domain.UserValue) (bool, error)
	ShowOneById(iId int) (*domain.User, error)
}
