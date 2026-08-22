package resource

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkg "example/pkg"
)

type GameTypeUsecase struct {
	*AbstractUsecase
	GameTypeModel outputPortAnyModel.GameTypeModel
}

func NewGameTypeUsecase(oGameTypeModel outputPortAnyModel.GameTypeModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyAdminResource.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameTypeModel:   oGameTypeModel,
	}
}

func (oSelf *GameTypeUsecase) AddOne(oValue *domain.GameTypeValue) (bool, error) {

	_, oErr := oSelf.GameTypeModel.AddOne(oValue)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) EditOne(oValue *domain.GameTypeValue, iId uint) (bool, error) {

	_, oErr := oSelf.GameTypeModel.EditOneById(oValue, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) RemoveOne(iId uint) (bool, error) {

	_, oErr := oSelf.GameTypeModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) ShowOne(iId uint) (*domain.GameType, error) {
	oGameType, oErr := oSelf.GameTypeModel.ShowOneById(iId)

	return oGameType, oErr
}

func (oSelf *GameTypeUsecase) ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, uint64, error) {

	aWheres := pkg.FiltersToMysqlWheres([]string{}, aFilters)
	aOrders := pkg.SortersToMysqlOrders([]string{}, aSorters)

	oLimit := pkg.PaginationToMysqlLimit(oPagination)

	aGameTypes, oErr := oSelf.GameTypeModel.ShowOnesByWheresWithOrdersLimit(aWheres, aOrders, oLimit)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.GameTypeModel.TotalByWheres(aWheres)

	return aGameTypes, iTotal, oErr
}
