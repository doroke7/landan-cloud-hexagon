package outputPortAnyLogic

import (
	domain "example/internal/domain"
)

type GameType interface {
	ShowTree() ([]*domain.GameType, uint64, error)
}
