package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pbResource "example/pb/resource"
)

func AdminPermissionGroup(oDomain *domain.AdminPermissionGroup) *pbResource.AdminPermissionGroup {
	if oDomain == nil {
		return nil
	}

	oProto := &pbResource.AdminPermissionGroup{
		Id:        oDomain.Id,
		ParentId:  oDomain.ParentId,
		Key:       oDomain.Key,
		Name:      oDomain.Name,
		CreatedAt: timestamppb.New(oDomain.CreatedAt),
		UpdatedAt: timestamppb.New(oDomain.UpdatedAt),
		DeletedAt: timestamppb.New(oDomain.DeletedAt),
		Parent:    AdminPermissionGroup(oDomain.Parent),
	}

	for i := range oDomain.Children {
		oChildPb := AdminPermissionGroup(&oDomain.Children[i])
		oProto.Children = append(oProto.Children, oChildPb)
	}

	for _, oAdminPermission := range oDomain.AdminPermissions {
		oAdminPermissionPb := adminPermission(oAdminPermission)
		oProto.AdminPermissions = append(oProto.AdminPermissions, oAdminPermissionPb)
	}

	return oProto
}
