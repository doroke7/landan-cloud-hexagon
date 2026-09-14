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

func (oSelf *GameTypeUsecase) AddOne(oAdd *domain.GameTypeVariable) error {

	oErr := oSelf.GameTypeModel.AddOne(oAdd)

	return oErr
}

func (oSelf *GameTypeUsecase) EditOneById(oEdit *domain.GameTypeVariable, iId uint64) error {

	oErr := oSelf.GameTypeModel.EditOneById(oEdit, iId)

	return oErr
}

func (oSelf *GameTypeUsecase) RemoveOneById(iId uint64) error {

	oErr := oSelf.GameTypeModel.RemoveOneById(iId)

	return oErr
}

func (oSelf *GameTypeUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.GameTypeModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
