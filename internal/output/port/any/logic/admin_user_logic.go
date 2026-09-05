package outputPortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminUserLogic interface {
	ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error)
	AddAminUser(oValue *domain.AdminUserValue) error
}
