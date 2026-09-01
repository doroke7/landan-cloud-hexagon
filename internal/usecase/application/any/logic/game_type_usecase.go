package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
)

type GameTypeUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.GameType
}

func NewGameTypeUsecase(oAbstractUsecase *AbstractUsecase, oGameTypeLogic outputPortAnyLogic.GameType) usecasePortAnyLogic.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameType:        oGameTypeLogic,
	}
}

func (oSelf *GameTypeUsecase) ShowTree() ([]*domain.GameType, error) {
	return oSelf.GameType.ShowTree()
}
