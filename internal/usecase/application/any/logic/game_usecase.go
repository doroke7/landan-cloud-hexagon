package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pkgInput "example/pkg/input"
)

type GameUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.GameLogic
}

func NewGameUsecase(oAbstractUsecase *AbstractUsecase, oGameLogic outputPortAnyLogic.GameLogic) usecasePortAnyLogic.GameUsecase {
	return &GameUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameLogic:       oGameLogic,
	}
}

func (oSelf *GameUsecase) ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error) {

	aGames, iTotal, oErr := oSelf.GameLogic.ShowGamesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGames, uint64(iTotal), oErr
}
