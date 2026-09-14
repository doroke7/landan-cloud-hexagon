package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameUsecase interface {
	AddOne(oValue *domain.GameVariable) error
	ShowOneByKey(sKey string) (*domain.Game, error)
	EditOneById(oValue *domain.GameVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
