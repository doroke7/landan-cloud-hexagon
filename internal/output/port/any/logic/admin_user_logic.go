package outputPortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminUserLogic interface {
	ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint, error)
	AddAminUser(oValue *domain.AdminUserValue) error
	EditAdminUserById(oValue *domain.AdminUserValue, iId uint) error
}
