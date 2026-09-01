package usecasePortAnyLogic

import (
	domain "example/internal/domain"
)

type GameTypeUsecase interface {
	ShowTree() ([]*domain.GameType, error)
}
