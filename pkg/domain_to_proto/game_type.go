package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pbResource "example/pb"
)

func GameType(oDomain *domain.GameType) *pbResource.GameType {
	if oDomain == nil {
		return nil
	}

	oProto := &pbResource.GameType{
		Id:        oDomain.Id,
		ParentId:  oDomain.ParentId,
		Key:       oDomain.Key,
		Name:      oDomain.Name,
		CreatedAt: timestamppb.New(oDomain.CreatedAt),
		UpdatedAt: timestamppb.New(oDomain.UpdatedAt),
		DeletedAt: timestamppb.New(oDomain.DeletedAt),
		Parent:    GameType(oDomain.Parent),
	}

	for i := range oDomain.Children {
		oProto.Children = append(oProto.Children, GameType(&oDomain.Children[i]))
	}

	return oProto
}
