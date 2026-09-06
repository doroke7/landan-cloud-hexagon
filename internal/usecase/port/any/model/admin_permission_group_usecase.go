package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupUsecase interface {
	AddOne(oAdminPermissionGroup *domain.AdminPermissionGroup) error
	ShowOneById(iId uint64) (*domain.AdminPermissionGroup, error)
	EditOneById(oAdminPermissionGroup *domain.AdminPermissionGroup, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, error)
}
