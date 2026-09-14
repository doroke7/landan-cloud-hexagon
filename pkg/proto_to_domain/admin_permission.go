package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb"
)

func AdminPermission(oProtoAdminPermission *pbResource.AdminPermission) domain.AdminPermission {
	if oProtoAdminPermission == nil {
		return domain.AdminPermission{}
	}

	oDomainAdminPermission := domain.AdminPermission{
		Id:                     uint64(oProtoAdminPermission.GetId()),
		AdminPermissionGroupId: uint64(oProtoAdminPermission.GetAdminPermissionGroupId()),
		Type:                   uint8(oProtoAdminPermission.GetType()),
		Key:                    oProtoAdminPermission.GetKey(),
		Name:                   oProtoAdminPermission.GetName(),
		CreatedAt:              oProtoAdminPermission.GetCreatedAt().AsTime(),
		UpdatedAt:              oProtoAdminPermission.GetUpdatedAt().AsTime(),
		DeletedAt:              oProtoAdminPermission.GetDeletedAt().AsTime(),
	}

	return oDomainAdminPermission
}
