package outputPortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminUserLogic interface {
	ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error)
	AddAminUser(oValue *domain.AdminUserVariable) error
	EditAdminUserById(oValue *domain.AdminUserVariable, iId uint64) error
	ShowAdminUserById(iId uint64) (*domain.AdminUser, error)
	RemoveAdminUserById(iId uint64) error
}
