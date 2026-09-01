package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
)

type GameTypeUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.GameTypeLogic
}

func NewGameTypeUsecase(oAbstractUsecase *AbstractUsecase, oGameTypeLogic outputPortAnyLogic.GameTypeLogic) usecasePortAnyLogic.GameTypeUsecase {
	return &GameTypeUsecase{
		AbstractUsecase: oAbstractUsecase,
		GameTypeLogic:   oGameTypeLogic,
	}
}

func (oSelf *GameTypeUsecase) ShowTree() ([]*domain.GameType, error) {
	return oSelf.GameTypeLogic.ShowTree()
}
