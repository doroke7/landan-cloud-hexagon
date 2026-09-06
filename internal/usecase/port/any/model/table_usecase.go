package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableUsecase interface {
	AddOne(oValue *domain.TableValue) error
	ShowOneById(iId uint64) (*domain.Table, error)
	EditOneById(oValue *domain.TableValue, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
