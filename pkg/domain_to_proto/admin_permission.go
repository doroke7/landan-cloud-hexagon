package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pbResource "example/pb/resource"
)

func AdminPermission(oDomainAdminPermission *domain.AdminPermission) *pbResource.AdminPermission {
	if oDomainAdminPermission == nil {
		return nil
	}

	oProtoAdminPermission := &pbResource.AdminPermission{
		Id:                     oDomainAdminPermission.Id,
		Type:                   uint64(oDomainAdminPermission.Type),
		Key:                    oDomainAdminPermission.Key,
		Name:                   oDomainAdminPermission.Name,
		CreatedAt:              timestamppb.New(oDomainAdminPermission.CreatedAt),
		UpdatedAt:              timestamppb.New(oDomainAdminPermission.UpdatedAt),
		DeletedAt:              timestamppb.New(oDomainAdminPermission.DeletedAt),
		AdminPermissionGroupId: oDomainAdminPermission.AdminPermissionGroupId,
	}

	return oProtoAdminPermission
}
