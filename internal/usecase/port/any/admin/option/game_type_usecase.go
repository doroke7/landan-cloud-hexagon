package any

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameTypeUsecase interface {
	ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, uint64, error)
}
