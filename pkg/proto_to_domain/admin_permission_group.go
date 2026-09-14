package protoToDomain

import (
	"example/internal/domain"

	pbResource "example/pb/resource"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func adminPermission(oNode *domain.AdminPermission) *pbResource.AdminPermission {
	if oNode == nil {
		return nil
	}

	oPb := &pbResource.AdminPermission{
		Id:        oNode.Id,
		Type:      uint64(oNode.Type),
		Key:       oNode.Key,
		Name:      oNode.Name,
		CreatedAt: timestamppb.New(oNode.CreatedAt),
		UpdatedAt: timestamppb.New(oNode.UpdatedAt),
		DeletedAt: timestamppb.New(oNode.DeletedAt),
	}

	return oPb
}

func AdminPermissionGroup(oNode *domain.AdminPermissionGroup) *pbResource.AdminPermissionGroup {
	if oNode == nil {
		return nil
	}

	oPb := &pbResource.AdminPermissionGroup{
		Id:        oNode.Id,
		ParentId:  oNode.ParentId,
		Key:       oNode.Key,
		Name:      oNode.Name,
		CreatedAt: timestamppb.New(oNode.CreatedAt),
		UpdatedAt: timestamppb.New(oNode.UpdatedAt),
		DeletedAt: timestamppb.New(oNode.DeletedAt),
		Parent:    AdminPermissionGroup(oNode.Parent),
	}

	for i := range oNode.Children {
		oChildPb := AdminPermissionGroup(&oNode.Children[i])
		oPb.Children = append(oPb.Children, oChildPb)
	}

	for _, oAdminPermission := range oNode.AdminPermissions {
		oAdminPermissionPb := adminPermission(oAdminPermission)
		oPb.AdminPermissions = append(oPb.AdminPermissions, oAdminPermissionPb)
	}

	return oPb
}
