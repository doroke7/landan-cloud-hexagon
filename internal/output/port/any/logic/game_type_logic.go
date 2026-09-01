package outputPortAnyLogic

import (
	domain "example/internal/domain"
)

type GameTypeLogic interface {
	ShowTree() ([]*domain.GameType, error)
}
