package option

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAdminOption "example/internal/usecase/port/any/admin/option"
	pkg "example/pkg"
)

type GameTypeUsecase struct {
	*AbstractUsecase
	GameTypeModel outputPortAnyModel.GameTypeModel
}

func NewGameTypeUsecase(oGameTypeModel outputPortAnyModel.GameTypeModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyAdminOption.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameTypeModel:   oGameTypeModel,
	}
}

func (oSelf *GameTypeUsecase) ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, uint64, error) {

	aWheres := pkg.FiltersToWheres([]string{}, aFilters)
	aOrders := pkg.SortersToOrders([]string{}, aSorters)

	oLimit := pkg.PaginationToLimit(oPagination)

	aGameTypes, oErr := oSelf.GameTypeModel.ShowOnesByWheresWithOrdersLimit(aWheres, aOrders, oLimit)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.GameTypeModel.TotalByWheres(aWheres)

	return aGameTypes, iTotal, oErr
}
