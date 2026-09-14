package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupUsecase interface {
	AddOne(oValue *domain.AdminPermissionGroupVariable) error
	ShowOne(iId uint64) (*domain.AdminPermissionGroup, error)
	ShowTree() ([]*domain.AdminPermissionGroup, error)
	EditOne(oValue *domain.AdminPermissionGroupVariable, iId uint64) error
	RemoveOne(iId uint64) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error)
}
