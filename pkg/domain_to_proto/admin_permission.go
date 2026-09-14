package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pb "example/pb"
)

func AdminPermission(oDomainAdminPermission *domain.AdminPermission) *pb.AdminPermission {
	if oDomainAdminPermission == nil {
		return nil
	}

	oProtoAdminPermission := &pb.AdminPermission{
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
