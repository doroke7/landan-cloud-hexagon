package outputPortAnyLogic

import (
	domain "example/internal/domain"
)

type AdminRoleLogic interface {
	AddAdminRole(oVariable *domain.AdminRoleVariable) error
	EditAdminRoleById(oVariable *domain.AdminRoleVariable, iId uint64) error
}
