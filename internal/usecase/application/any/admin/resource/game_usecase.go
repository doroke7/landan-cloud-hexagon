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
			return errors.New("key already used by another game")
		}
	}

	oErr := oSelf.GameModel.AddOne(oValue)

	return oErr
}

func (oSelf *GameUsecase) EditOne(oValue *domain.GameValue, iId uint64) error {

	if oValue.Key != nil {
		oGameByKey, oErr := oSelf.GameModel.ShowOneByKey(*oValue.Key)

		if oErr != nil {
			return oErr
		}

		if oGameByKey != nil && oGameByKey.Id != iId {
			return pkgUtility.NewDefaultError("key="+*oValue.Key+" already exists", -2, 500)
		}
	}

	oErr := oSelf.GameModel.EditOneById(oValue, iId)

	return oErr
}

func (oSelf *GameUsecase) RemoveOne(iId uint) error {

	oErr := oSelf.GameModel.RemoveOneById(uint64(iId))

	return oErr
}

func (oSelf *GameUsecase) ShowOne(iId uint64) (*domain.Game, error) {
	oGame, oErr := oSelf.GameModel.ShowOneById(iId)

	return oGame, oErr
}

func (oSelf *GameUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error) {

	aGames, iTotal, oErr := oSelf.GameLogic.ShowGamesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGames, uint64(iTotal), oErr
}
