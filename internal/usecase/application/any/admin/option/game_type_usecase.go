package usecaseApplicationAnyAdminOption

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminOption "example/internal/usecase/port/any/admin/option"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	GameTypeModel outputPortAnyModel.GameTypeModel
}

func NewGameTypeUsecase(oGameTypeModel outputPortAnyModel.GameTypeModel, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminOption.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameTypeModel:   oGameTypeModel,
	}
}

func (oSelf *GameTypeUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error) {

	aGameTypes, oErr := oSelf.GameTypeModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.GameTypeModel.TotalByFilters(aFilters)

	return aGameTypes, iTotal, oErr
}
