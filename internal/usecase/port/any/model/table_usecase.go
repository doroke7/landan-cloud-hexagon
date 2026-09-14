package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableUsecase interface {
	AddOne(oValue *domain.TableVariable) error
	EditOneById(oValue *domain.TableVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
