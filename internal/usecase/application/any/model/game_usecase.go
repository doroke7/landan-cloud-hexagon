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

func (oSelf *GameUsecase) AddOne(oAdd *domain.GameValue) error {

	oErr := oSelf.GameModel.AddOne(oAdd)

	return oErr
}

func (oSelf *GameUsecase) ShowOneById(iId uint64) (*domain.Game, error) {

	oGame, oErr := oSelf.GameModel.ShowOneById(iId)

	return oGame, oErr
}

func (oSelf *GameUsecase) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error) {

	aGames, oErr := oSelf.GameModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGames, oErr
}

func (oSelf *GameUsecase) EditOneById(oEdit *domain.GameValue, iId uint64) error {

	oErr := oSelf.GameModel.EditOneById(oEdit, iId)

	return oErr
}

func (oSelf *GameUsecase) RemoveOneById(iId uint64) error {

	oErr := oSelf.GameModel.RemoveOneById(iId)

	return oErr
}

func (oSelf *GameUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.GameModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
