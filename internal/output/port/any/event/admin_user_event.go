package outputPortAnyEvent

import (
	domain "example/internal/domain"
)

type AdminUserEvent interface {
	AddOne(oAdminUser *domain.AdminUser) error
}
