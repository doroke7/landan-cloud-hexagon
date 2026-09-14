package domainToProto

import (
	domain "example/internal/domain"
	pbResourceModel "example/pb/resource/model"
)

func AdminPermissionValue(oDomainAdminPermissionValue *domain.AdminPermissionValue) *pbResourceModel.AdminPermissionValue {
	oProtoAdminPermissionValue := &pbResourceModel.AdminPermissionValue{
		Key:  oDomainAdminPermissionValue.Key,
		Name: oDomainAdminPermissionValue.Name,
	}

	if oDomainAdminPermissionValue.Type != nil {
		iType := uint64(*oDomainAdminPermissionValue.Type)
		oProtoAdminPermissionValue.Type = &iType
	}

	return oProtoAdminPermissionValue
}
