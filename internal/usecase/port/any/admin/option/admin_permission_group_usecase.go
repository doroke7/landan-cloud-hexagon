package usecasePortAnyAdminOption

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupUsecase interface {
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error)
	ShowTree() ([]*domain.AdminPermissionGroup, error)
}
