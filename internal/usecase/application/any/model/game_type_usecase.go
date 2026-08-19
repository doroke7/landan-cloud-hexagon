package usecase

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkg "example/pkg"
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

	_, oErr := oSelf.GameTypeModel.AddOne(oAdd)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) ShowOneById(iId uint) (*domain.GameType, error) {

	oGameType, oErr := oSelf.GameTypeModel.ShowOneById(iId)

	return oGameType, oErr
}

func (oSelf *GameTypeUsecase) ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.GameType, error) {

	aGameTypes, oErr := oSelf.GameTypeModel.ShowOnesByWheresWithOrdersLimit(aWheres, aOrders, oLimit)

	return aGameTypes, oErr
}

func (oSelf *GameTypeUsecase) EditOneById(oEdit *domain.GameTypeValue, iId uint64) (bool, error) {

	_, oErr := oSelf.GameTypeModel.EditOneById(oEdit, uint(iId))

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) RemoveOneById(iId uint) (bool, error) {

	_, oErr := oSelf.GameTypeModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) TotalByWheres(aWheres []*pkg.Where) (uint64, error) {

	iTotal, oErr := oSelf.GameTypeModel.TotalByWheres(aWheres)

	return iTotal, oErr
}
