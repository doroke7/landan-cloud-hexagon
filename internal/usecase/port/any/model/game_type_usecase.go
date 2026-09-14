package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase interface {
	AddOne(oValue *domain.GameTypeVariable) error
	EditOneById(oValue *domain.GameTypeVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
