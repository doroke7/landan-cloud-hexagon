package resource

import (
	domain "example/internal/domain"
	outputPortAnylogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"
	"fmt"

	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkg "example/pkg"
)

type GameUsecase struct {
	*AbstractUsecase
	GameModel outputPortAnyModel.GameModel
	GameLogic outputPortAnylogic.GameLogic
}

func NewGameUsecase(oGameModel outputPortAnyModel.GameModel, oGameLogic outputPortAnylogic.GameLogic, oAbstractUsecase *AbstractUsecase) usecasePortAnyAdminResource.GameUsecase {
	return &GameUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameModel:       oGameModel,
		GameLogic:       oGameLogic,
	}
}

func (oSelf *GameUsecase) AddOne(oValue *domain.GameValue) (bool, error) {
	fmt.Println("25....")
	_, oErr := oSelf.GameModel.AddOne(oValue)
	fmt.Println("27....")

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) EditOne(oValue *domain.GameValue, iId uint) (bool, error) {

	_, oErr := oSelf.GameModel.EditOneById(oValue, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) RemoveOne(iId uint) (bool, error) {

	_, oErr := oSelf.GameModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) ShowOne(iId uint) (*domain.Game, error) {
	oGame, oErr := oSelf.GameModel.ShowOneById(iId)

	return oGame, oErr
}

func (oSelf *GameUsecase) ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Game, uint64, error) {

	aWheres := pkg.FiltersToMysqlWheres([]string{}, aFilters)
	aOrders := pkg.SortersToMysqlOrders([]string{}, aSorters)

	oGameLimit := pkg.PaginationToMysqlLimit(oPagination)

	aGames, iTotal, oErr := oSelf.GameLogic.ShowGamesTotalByWheresWithOrdersLimit(aWheres, aOrders, oGameLimit)

	return aGames, uint64(iTotal), oErr
}
