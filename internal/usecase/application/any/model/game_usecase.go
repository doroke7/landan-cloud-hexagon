package usecase

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkg "example/pkg"
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

	_, oErr := oSelf.GameModel.AddOne(oAdd)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) ShowOneById(iId uint) (*domain.Game, error) {

	oGame, oErr := oSelf.GameModel.ShowOneById(iId)

	return oGame, oErr
}

func (oSelf *GameUsecase) ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Game, error) {

	aGames, oErr := oSelf.GameModel.ShowOnesByWheresWithOrdersLimit(aWheres, aOrders, oLimit)

	return aGames, oErr
}

func (oSelf *GameUsecase) EditOneById(oEdit *domain.GameValue, iId uint64) (bool, error) {

	_, oErr := oSelf.GameModel.EditOneById(oEdit, uint(iId))

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) RemoveOneById(iId uint) (bool, error) {

	_, oErr := oSelf.GameModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) TotalByWheres(aWheres []*pkg.Where) (uint64, error) {

	iTotal, oErr := oSelf.GameModel.TotalByWheres(aWheres)

	return iTotal, oErr
}
