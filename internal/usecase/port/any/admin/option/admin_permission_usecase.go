package usecasePortAnyAdminOption

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionUsecase interface {
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, uint64, error)
}
