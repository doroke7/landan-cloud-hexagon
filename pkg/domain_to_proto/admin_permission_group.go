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

	for _, oChild := range oDomain.Children {
		oChildPb := AdminPermissionGroup(oChild)
		oProto.Children = append(oProto.Children, oChildPb)
	}

	for _, oAdminPermission := range oDomain.AdminPermissions {
		oAdminPermissionPb := AdminPermission(oAdminPermission)
		oProto.AdminPermissions = append(oProto.AdminPermissions, oAdminPermissionPb)
	}

	return oProto
}
