package usecaseApplicationAnyAdminOption

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminOption "example/internal/usecase/port/any/admin/option"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	GameTypeModel outputPortAnyModel.GameTypeModel
	GameTypeLogic outputPortAnyLogic.GameTypeLogic
}

func NewGameTypeUsecase(oGameTypeModel outputPortAnyModel.GameTypeModel, oGameTypeLogic outputPortAnyLogic.GameTypeLogic, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminOption.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameTypeModel:   oGameTypeModel,
		GameTypeLogic:   oGameTypeLogic,
	}
}

func (oSelf *GameTypeUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error) {

	aGameTypes, oErr := oSelf.GameTypeLogic.ShowGameTypes()
	if oErr != nil {
		return nil, 0, oErr
	}

	return aGameTypes, uint64(0), oErr
}

func (oSelf *GameTypeUsecase) ShowTree() ([]*domain.GameType, error) {
	aTree, oErr := oSelf.GameTypeLogic.ShowTree()

	return aTree, oErr
}
