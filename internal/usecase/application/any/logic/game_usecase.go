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

func (oSelf *GameUsecase) ShowGamesTotalByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Game, int64, error) {

	aGames, iTotal, oErr := oSelf.GameLogic.ShowGamesTotalByWheresWithOrdersLimit(aWheres, aOrders, oLimit)

	return aGames, iTotal, oErr
}
