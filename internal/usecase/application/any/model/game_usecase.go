package usecaseApplicationAnyModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type GameUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.GameModel
}

func NewGameUsecase(oGameModel outputPortAnyModel.GameModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyModel.GameUsecase {
	return &GameUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameModel:       oGameModel,
	}
}

func (oSelf *GameUsecase) AddOne(oAdd *domain.GameValue) (bool, error) {

	oErr := oSelf.GameModel.AddOne(oAdd)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) ShowOneById(iId uint) (*domain.Game, error) {

	oGame, oErr := oSelf.GameModel.ShowOneById(iId)

	return oGame, oErr
}

func (oSelf *GameUsecase) ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error) {

	aGames, oErr := oSelf.GameModel.ShowOnesByFiltersWithOrdersPagination(aFilters, aSorters, oPagination)

	return aGames, oErr
}

func (oSelf *GameUsecase) EditOneById(oEdit *domain.GameValue, iId uint64) (bool, error) {

	oErr := oSelf.GameModel.EditOneById(oEdit, uint(iId))

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) RemoveOneById(iId uint) (bool, error) {

	oErr := oSelf.GameModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.GameModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
