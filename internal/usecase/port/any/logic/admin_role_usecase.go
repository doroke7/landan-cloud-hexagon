package usecasePortAnyLogic

import (
	domain "example/internal/domain"
)

type AdminRoleUsecase interface {
	AddAdminRole(oVariable *domain.AdminRoleVariable) error
}
