package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableUsecase interface {
	AddOne(oValue *domain.TableValue) (bool, error)
	ShowOne(iId uint) (*domain.Table, error)
	EditOne(oValue *domain.TableValue, iId uint) (bool, error)
	RemoveOne(iId uint) (bool, error)
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, uint64, error)
}
