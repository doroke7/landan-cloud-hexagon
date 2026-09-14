package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb"
)

func Game(oProtoGame *pbResource.Game) domain.Game {
	if oProtoGame == nil {
		return domain.Game{}
	}

	return domain.Game{
		Id:          uint64(oProtoGame.GetId()),
		GameTypeId:  uint64(oProtoGame.GetGameTypeId()),
		Key:         oProtoGame.GetKey(),
		Name:        oProtoGame.GetName(),
		Description: oProtoGame.GetDescription(),
		CreatedAt:   oProtoGame.GetCreatedAt().AsTime(),
		UpdatedAt:   oProtoGame.GetUpdatedAt().AsTime(),
		DeletedAt:   oProtoGame.GetDeletedAt().AsTime(),
		GameType:    GameType(oProtoGame.GetGameType()),
	}
}
