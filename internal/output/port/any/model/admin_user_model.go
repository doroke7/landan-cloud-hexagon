package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminUserModel interface {
	ShowOneByName(sName string) (*domain.AdminUser, error)
	ShowOneById(iId uint) (*domain.AdminUser, error)
	AddOne(oAdminUser *domain.AdminUserValue) (bool, error)
	EditOneById(oAdminUser *domain.AdminUserValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
