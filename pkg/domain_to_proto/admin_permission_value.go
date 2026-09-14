package domainToProto

import (
	domain "example/internal/domain"
	pbResource "example/pb/resource"
)

func AdminPermissionValue(oDomainAdminPermissionValue *domain.AdminPermissionValue) *pbResource.AdminPermissionValue {
	oProtoAdminPermissionValue := &pbResource.AdminPermissionValue{
		Key:  oDomainAdminPermissionValue.Key,
		Name: oDomainAdminPermissionValue.Name,
	}

	if oDomainAdminPermissionValue.Type != nil {
		iType := uint64(*oDomainAdminPermissionValue.Type)
		oProtoAdminPermissionValue.Type = &iType
	}

	return oProtoAdminPermissionValue
}
