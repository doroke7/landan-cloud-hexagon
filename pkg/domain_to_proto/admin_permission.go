package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pbResource "example/pb/resource"
)

func AdminPermission(oDomain *domain.AdminPermission) *pbResource.AdminPermission {
	if oDomain == nil {
		return nil
	}

	oProto := &pbResource.AdminPermission{
		Id:                     oDomain.Id,
		Type:                   uint64(oDomain.Type),
		Key:                    oDomain.Key,
		Name:                   oDomain.Name,
		CreatedAt:              timestamppb.New(oDomain.CreatedAt),
		UpdatedAt:              timestamppb.New(oDomain.UpdatedAt),
		DeletedAt:              timestamppb.New(oDomain.DeletedAt),
		AdminPermissionGroupId: oDomain.AdminPermissionGroupId,
	}

	return oProto
}
