package usecase

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pkg "example/pkg"
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

func (oSelf *GameUsecase) ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Game, int64, error) {

	aGames, iTotal, oErr := oSelf.GameLogic.ShowGamesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGames, iTotal, oErr
}
