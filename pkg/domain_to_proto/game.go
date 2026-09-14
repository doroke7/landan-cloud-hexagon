package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pb "example/pb"
)

func Game(oDomain *domain.Game) *pb.Game {
	if oDomain == nil {
		return nil
	}

	return &pb.Game{
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
