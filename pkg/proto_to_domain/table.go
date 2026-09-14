package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb/resource"
)

func Table(oProtoTable *pbResource.Table) domain.Table {
	if oProtoTable == nil {
		return domain.Table{}
	}

	oGame := Game(oProtoTable.GetGame())
	oTable := domain.Table{
		Id:          uint64(oProtoTable.GetId()),
		No:          oProtoTable.GetNo(),
		GameId:      uint64(oProtoTable.GetGameId()),
		Key:         oProtoTable.GetKey(),
		State:       uint8(oProtoTable.GetState()),
		Description: oProtoTable.GetDescription(),
		Result:      oProtoTable.GetResult(),
		StartedAt:   oProtoTable.GetStartedAt().AsTime(),
		EndedAt:     oProtoTable.GetEndedAt().AsTime(),
		CreatedAt:   oProtoTable.GetCreatedAt().AsTime(),
		UpdatedAt:   oProtoTable.GetUpdatedAt().AsTime(),
		DeletedAt:   oProtoTable.GetDeletedAt().AsTime(),
		Game:        oGame,
	}

	return oTable
}
