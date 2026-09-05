package usecaseApplicationAnyAdminResource

import (
	domain "example/internal/domain"
	outputPortAnylogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameTypeUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	GameTypeModel outputPortAnyModel.GameTypeModel
	GameTypeLogic outputPortAnylogic.GameTypeLogic
	GameModel     outputPortAnyModel.GameModel
}

func NewGameTypeUsecase(oGameTypeModel outputPortAnyModel.GameTypeModel, oGameTypeLogic outputPortAnylogic.GameTypeLogic, oGameModel outputPortAnyModel.GameModel, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameTypeModel:   oGameTypeModel,
		GameTypeLogic:   oGameTypeLogic,
		GameModel:       oGameModel,
	}
}

func (oSelf *GameTypeUsecase) AddOne(oValue *domain.GameTypeValue) (bool, error) {

	oErr := oSelf.GameTypeModel.AddOne(oValue)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) EditOne(oValue *domain.GameTypeValue, iId uint) (bool, error) {

	oErr := oSelf.GameTypeModel.EditOneById(oValue, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) ShowOne(iId uint) (*domain.GameType, error) {
	oGameType, oErr := oSelf.GameTypeModel.ShowOneById(iId)

	return oGameType, oErr
}

// ShowTree 回一個虛擬 root，Children 掛所有頂層 game type（整棵樹）。
func (oSelf *GameTypeUsecase) ShowTree() (*domain.GameType, error) {

	aRoots, oErr := oSelf.GameTypeLogic.ShowTree()
	if oErr != nil {
		return nil, oErr
	}

	oRoot := &domain.GameType{}
	for _, oOne := range aRoots {
		oRoot.Children = append(oRoot.Children, *oOne)
	}

	return oRoot, nil
}

func (oSelf *GameTypeUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error) {

	aGameTypes, iTotal, oErr := oSelf.GameTypeLogic.ShowGameTypesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGameTypes, uint64(iTotal), oErr
}

func (oSelf *GameTypeUsecase) RemoveOne(iId uint) (bool, error) {

	// 底下還有未刪除的子類型就不放行，避免刪掉父類型後留下孤兒資料
	iChildren, oErr := oSelf.GameTypeModel.TotalByParentId(iId)

	if oErr != nil {
		return false, oErr
	}

	if iChildren > 0 {
		return false, pkgUtility.NewDefaultError("This game type still has child types; remove them before deleting", -2, 500)
	}

	// 底下還有掛在此類型的遊戲也不放行
	iGames, oErr := oSelf.GameModel.TotalByGameTypeId(iId)

	if oErr != nil {
		return false, oErr
	}

	if iGames > 0 {
		return false, pkgUtility.NewDefaultError("This game type still has games; remove them before deleting", -2, 500)
	}

	oErr = oSelf.GameTypeModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}
