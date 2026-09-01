package usecasePortAnyEvent

import (
	domain "example/internal/domain"
)

type AdminUserUsecase interface {
	AddOne(oAdminUser *domain.AdminUser) error
}
