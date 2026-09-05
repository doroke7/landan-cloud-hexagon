package usecaseApplicationAnyModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.GameTypeModel
}

func NewGameTypeUsecase(oGameTypeModel outputPortAnyModel.GameTypeModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyModel.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameTypeModel:   oGameTypeModel,
	}
}

func (oSelf *GameTypeUsecase) AddOne(oAdd *domain.GameTypeValue) (bool, error) {

	oErr := oSelf.GameTypeModel.AddOne(oAdd)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) ShowOneById(iId uint) (*domain.GameType, error) {

	oGameType, oErr := oSelf.GameTypeModel.ShowOneById(iId)

	return oGameType, oErr
}

func (oSelf *GameTypeUsecase) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error) {

	aGameTypes, oErr := oSelf.GameTypeModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGameTypes, oErr
}

func (oSelf *GameTypeUsecase) EditOneById(oEdit *domain.GameTypeValue, iId uint64) (bool, error) {

	oErr := oSelf.GameTypeModel.EditOneById(oEdit, uint(iId))

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) RemoveOneById(iId uint) (bool, error) {

	oErr := oSelf.GameTypeModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.GameTypeModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
