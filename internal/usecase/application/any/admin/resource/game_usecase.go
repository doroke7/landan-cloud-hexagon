package usecaseApplicationAnyAdminResource

import (
	"errors"

	domain "example/internal/domain"
	outputPortAnylogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"

	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	GameModel outputPortAnyModel.GameModel
	GameLogic outputPortAnylogic.GameLogic
}

func NewGameUsecase(oGameModel outputPortAnyModel.GameModel, oGameLogic outputPortAnylogic.GameLogic, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.GameUsecase {
	return &GameUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameModel:       oGameModel,
		GameLogic:       oGameLogic,
	}
}

func (oSelf *GameUsecase) AddOne(oValue *domain.GameValue) (bool, error) {

	if oValue.Key != nil {
		oGameByKey, oErr := oSelf.GameModel.ShowOneByKey(*oValue.Key)

		if oErr != nil {
			return false, oErr
		}

		if oGameByKey != nil {
			return false, errors.New("key 已被其他遊戲使用")
		}
	}

	_, oErr := oSelf.GameModel.AddOne(oValue)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameUsecase) EditOne(oValue *domain.GameValue, iId uint) (bool, error) {

	if oValue.Key != nil {
		oGameByKey, oErr := oSelf.GameModel.ShowOneByKey(*oValue.Key)

		if oErr != nil {
			return false, oErr
		}

		if oGameByKey != nil && oGameByKey.Id != iId {
			return false, pkgUtility.NewDefaultError("key="+*oValue.Key+" 已經存在", -4, 500)
		}
	}

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

func (oSelf *GameUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error) {

	aGames, iTotal, oErr := oSelf.GameLogic.ShowGamesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGames, uint64(iTotal), oErr
}
