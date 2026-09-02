package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameUsecase interface {
	AddOne(oValue *domain.GameValue) (bool, error)
	ShowOneById(iId uint) (*domain.Game, error)
	ShowOneByKey(sKey string) (*domain.Game, error)
	EditOneById(oValue *domain.GameValue, iId uint64) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
	ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error)
}
