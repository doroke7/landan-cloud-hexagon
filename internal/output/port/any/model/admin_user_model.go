package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type AdminUserModel interface {
	ShowOneByName(sName string) (*domain.AdminUser, error)
	ShowOneById(iId uint) (*domain.AdminUser, error)
	AddOne(oAdminUser *domain.AdminUserValue) (bool, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.AdminUser, error)
	TotalByFilters(aFilters []*pkg.Filter) (uint64, error)
}
