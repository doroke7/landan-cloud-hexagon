package any

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameTypeUsecase interface {
	AddOne(oValue *domain.GameTypeValue) (bool, error)
	ShowOne(iId uint) (*domain.GameType, error)
	EditOne(oValue *domain.GameTypeValue, iId uint) (bool, error)
	RemoveOne(iId uint) (bool, error)
	ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, uint64, error)
}
