package domainToProto

import (
	domain "example/internal/domain"
	pbResource "example/pb/resource"
)

func AdminPermissionValue(oDomainAdminPermissionValue *domain.AdminPermissionVariable) *pbResource.AdminPermissionValue {
	if oDomainAdminPermissionValue == nil {
		return nil
	}

	oProtoAdminPermissionValue := &pbResource.AdminPermissionValue{
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
