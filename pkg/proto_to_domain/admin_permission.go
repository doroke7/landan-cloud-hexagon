package protoToDomain

import (
	"example/internal/domain"

	pbResource "example/pb/resource"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func adminPermission(oDomain *domain.AdminPermission) *pbResource.AdminPermission {
	if oDomain == nil {
		return nil
	}

	oProto := &pbResource.AdminPermission{
		Id:        oDomain.Id,
		Type:      uint64(oDomain.Type),
		Key:       oDomain.Key,
		Name:      oDomain.Name,
		CreatedAt: timestamppb.New(oDomain.CreatedAt),
		UpdatedAt: timestamppb.New(oDomain.UpdatedAt),
		DeletedAt: timestamppb.New(oDomain.DeletedAt),
	}

	return oProto
}
