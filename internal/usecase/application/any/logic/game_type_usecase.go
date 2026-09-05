package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.GameTypeLogic
}

func NewGameTypeUsecase(oAbstractUsecase *AbstractUsecase, oGameTypeLogic outputPortAnyLogic.GameTypeLogic) usecasePortAnyLogic.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameTypeLogic:   oGameTypeLogic,
	}
}

func (oSelf *GameTypeUsecase) ShowTree() ([]*domain.GameType, error) {
	return oSelf.GameTypeLogic.ShowTree()
}

func (oSelf *GameTypeUsecase) ShowGameTypesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error) {

	aGameTypes, iTotal, oErr := oSelf.GameTypeLogic.ShowGameTypesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGameTypes, uint64(iTotal), oErr
}
