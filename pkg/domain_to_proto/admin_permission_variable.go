package domainToProto

import (
	domain "example/internal/domain"
	pbResource "example/pb/resource"
)

func AdminPermissionVariable(oDomainAdminPermissionValue *domain.AdminPermissionVariable) *pbResource.AdminPermissionVariable {
	if oDomainAdminPermissionValue == nil {
		return nil
	}

	oProtoAdminPermissionValue := &pbResource.AdminPermissionVariable{
		Id:   oDomainAdminPermissionValue.Id,
		Key:  oDomainAdminPermissionValue.Key,
		Name: oDomainAdminPermissionValue.Name,
	}

	if oDomainAdminPermissionValue.Type != nil {
		iType := uint64(*oDomainAdminPermissionValue.Type)
		oProtoAdminPermissionValue.Type = &iType
	}

	return oProtoAdminPermissionValue
}
