package usecasePortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase interface {
	ShowTree() ([]*domain.GameType, error)
	ShowGameTypeById(iId uint64) (*domain.GameType, error)
	ShowGameTypesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error)
}
