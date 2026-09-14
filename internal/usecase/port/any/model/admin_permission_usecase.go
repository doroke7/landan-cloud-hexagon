package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionUsecase interface {
	AddOne(oAdminPermission *domain.AdminPermissionVariable) error
	ShowOneById(iId uint64) (*domain.AdminPermission, error)
	EditOneById(oAdminPermission *domain.AdminPermissionVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
