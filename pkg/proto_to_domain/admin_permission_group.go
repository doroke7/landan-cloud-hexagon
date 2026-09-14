package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb"
)

func AdminPermissionGroup(oProto *pbResource.AdminPermissionGroup) domain.AdminPermissionGroup {
	if oProto == nil {
		return domain.AdminPermissionGroup{}
	}

	oDomainAdminPermissionGroup := domain.AdminPermissionGroup{
		Id:        uint64(oProto.GetId()),
		ParentId:  uint64(oProto.GetParentId()),
		Key:       oProto.GetKey(),
		Name:      oProto.GetName(),
		CreatedAt: oProto.GetCreatedAt().AsTime(),
		UpdatedAt: oProto.GetUpdatedAt().AsTime(),
		DeletedAt: oProto.GetDeletedAt().AsTime(),
		// 沒有子節點時也回 []（非 nil），避免 JSON 出現 children: null
		Children: make([]*domain.AdminPermissionGroup, 0, len(oProto.GetChildren())),
	}

	if oProtoParent := oProto.GetParent(); oProtoParent != nil {
		oDomainParent := AdminPermissionGroup(oProtoParent)
		oDomainAdminPermissionGroup.Parent = &oDomainParent
	}

	for _, oProtoChild := range oProto.GetChildren() {
		oDomainChild := AdminPermissionGroup(oProtoChild)
		oDomainAdminPermissionGroup.Children = append(oDomainAdminPermissionGroup.Children, &oDomainChild)
	}

	for _, oProtoAdminPermission := range oProto.GetAdminPermissions() {
		oDomainAdminPermission := AdminPermission(oProtoAdminPermission)
		oDomainAdminPermissionGroup.AdminPermissions = append(oDomainAdminPermissionGroup.AdminPermissions, &oDomainAdminPermission)
	}

	return oDomainAdminPermissionGroup
}
