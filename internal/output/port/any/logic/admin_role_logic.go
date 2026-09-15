package outputPortAnyLogic

import (
	domain "example/internal/domain"
)

type AdminRoleLogic interface {
	AddAdminRole(oVariable *domain.AdminRoleVariable) error
}
