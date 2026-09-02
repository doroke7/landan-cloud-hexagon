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

	_, oErr := oSelf.GameTypeModel.AddOne(oValue)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) EditOne(oValue *domain.GameTypeValue, iId uint) (bool, error) {

	_, oErr := oSelf.GameTypeModel.EditOneById(oValue, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeUsecase) ShowOne(iId uint) (*domain.GameType, error) {
	oGameType, oErr := oSelf.GameTypeModel.ShowOneById(iId)

	return oGameType, oErr
}

func (oSelf *GameTypeUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error) {

	aGameTypes, iTotal, oErr := oSelf.GameTypeLogic.ShowGameTypesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aGameTypes, iTotal, oErr
}

func (oSelf *GameTypeUsecase) RemoveOne(iId uint) (bool, error) {

	// 底下還有未刪除的子類型就不放行，避免刪掉父類型後留下孤兒資料
	iChildren, oErr := oSelf.GameTypeModel.TotalByParentId(iId)

	if oErr != nil {
		return false, oErr
	}

	if iChildren > 0 {
		return false, pkgUtility.NewDefaultError("此遊戲類型底下還有子類型，請先移除子類型再刪除", -2, 500)
	}

	// 底下還有掛在此類型的遊戲也不放行
	iGames, oErr := oSelf.GameModel.TotalByGameTypeId(iId)

	if oErr != nil {
		return false, oErr
	}

	if iGames > 0 {
		return false, pkgUtility.NewDefaultError("此遊戲類型底下還有遊戲，請先移除遊戲再刪除", -2, 500)
	}

	_, oErr = oSelf.GameTypeModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}
