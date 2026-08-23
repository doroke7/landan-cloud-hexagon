package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type TableUsecase interface {
	AddOne(oValue *domain.TableValue) (bool, error)
	ShowOne(iId uint) (*domain.Table, error)
	EditOne(oValue *domain.TableValue, iId uint) (bool, error)
	RemoveOne(iId uint) (bool, error)
	ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, uint64, error)
}
