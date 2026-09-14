package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pbResource "example/pb/resource"
)

func Game(oDomain *domain.Game) *pbResource.Game {
	if oDomain == nil {
		return nil
	}

	return &pbResource.Game{
		Id:          oDomain.Id,
		GameTypeId:  oDomain.GameTypeId,
		Key:         oDomain.Key,
		Name:        oDomain.Name,
		Description: oDomain.Description,
		CreatedAt:   timestamppb.New(oDomain.CreatedAt),
		UpdatedAt:   timestamppb.New(oDomain.UpdatedAt),
		DeletedAt:   timestamppb.New(oDomain.DeletedAt),
		GameType:    GameType(&oDomain.GameType),
	}
}
