package usecaseApplicationAnyAdminResource

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameTypeUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	GameTypeModel outputPortAnyModel.GameTypeModel
}

func NewGameTypeUsecase(oGameTypeModel outputPortAnyModel.GameTypeModel, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameTypeModel:   oGameTypeModel,
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

	aGameTypes, oErr := oSelf.GameTypeModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.GameTypeModel.TotalByFilters(aFilters)

	return aGameTypes, iTotal, oErr
}

func (oSelf *GameTypeUsecase) RemoveOne(iId uint) (bool, error) {

	// 底下還有未刪除的子類型就不放行，避免刪掉父類型後留下孤兒資料
	aChildren, oErr := oSelf.GameTypeModel.ShowOnesByParentId(iId)

	if oErr != nil {
		return false, oErr
	}

	if len(aChildren) > 0 {
		return false, pkgUtility.NewDefaultError("此遊戲類型底下還有子類型，請先移除子類型再刪除", -4, 500)
	}

	_, oErr = oSelf.GameTypeModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}
