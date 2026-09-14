package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb/resource"
)

func AdminPermissionGroup(oProto *pbResource.AdminPermissionGroup) domain.AdminPermissionGroup {
	if oProto == nil {
		return domain.AdminPermissionGroup{}
	}

	oAdminPermissionGroup := domain.AdminPermissionGroup{
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

	if oParent := oProto.GetParent(); oParent != nil {
		oParentDomain := AdminPermissionGroup(oParent)
		oAdminPermissionGroup.Parent = &oParentDomain
	}

	for _, oChild := range oProto.GetChildren() {
		oChildDomain := AdminPermissionGroup(oChild)
		oAdminPermissionGroup.Children = append(oAdminPermissionGroup.Children, &oChildDomain)
	}

	return oAdminPermissionGroup
}
