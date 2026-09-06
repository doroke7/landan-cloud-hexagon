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

	aTree, oErr := oSelf.GameTypeLogic.ShowTree()

	return aTree, oErr
}

func (oSelf *GameTypeUsecase) ShowGameTypeById(iId uint64) (*domain.GameType, error) {

	oGameType, oErr := oSelf.GameTypeLogic.ShowGameTypeById(iId)

	return oGameType, oErr
}

func (oSelf *GameTypeUsecase) ShowGameTypesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error) {

	aGameTypes, iTotal, oErr := oSelf.GameTypeLogic.ShowGameTypesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGameTypes, uint64(iTotal), oErr
}
