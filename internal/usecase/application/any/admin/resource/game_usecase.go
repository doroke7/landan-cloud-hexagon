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

func (oSelf *GameUsecase) AddOne(oValue *domain.GameValue) error {

	if oValue.Key != nil {
		oGameByKey, oErr := oSelf.GameModel.ShowOneByKey(*oValue.Key)

		if oErr != nil {
			return oErr
		}

		if oGameByKey != nil {
			return errors.New("key 已被其他遊戲使用")
		}
	}

	return oSelf.GameModel.AddOne(oValue)
}

func (oSelf *GameUsecase) EditOne(oValue *domain.GameValue, iId uint) error {

	if oValue.Key != nil {
		oGameByKey, oErr := oSelf.GameModel.ShowOneByKey(*oValue.Key)

		if oErr != nil {
			return oErr
		}

		if oGameByKey != nil && oGameByKey.Id != iId {
			return pkgUtility.NewDefaultError("key="+*oValue.Key+" already exists", -2, 500)
		}
	}

	return oSelf.GameModel.EditOneById(oValue, iId)
}

func (oSelf *GameUsecase) RemoveOne(iId uint) error {

	return oSelf.GameModel.RemoveOneById(iId)
}

func (oSelf *GameUsecase) ShowOne(iId uint) (*domain.Game, error) {
	oGame, oErr := oSelf.GameModel.ShowOneById(iId)

	return oGame, oErr
}

func (oSelf *GameUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint, error) {

	aGames, iTotal, oErr := oSelf.GameLogic.ShowGamesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGames, iTotal, oErr
}
