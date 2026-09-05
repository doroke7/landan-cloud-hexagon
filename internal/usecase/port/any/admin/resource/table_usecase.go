package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableUsecase interface {
	AddOne(oValue *domain.TableValue) error
	ShowOne(iId uint) (*domain.Table, error)
	EditOne(oValue *domain.TableValue, iId uint) error
	RemoveOne(iId uint) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, uint, error)
}
