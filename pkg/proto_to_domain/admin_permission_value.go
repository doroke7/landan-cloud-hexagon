package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb/resource"
)

func AdminPermissionValue(oProto *pbResource.AdminPermissionValue) domain.AdminPermissionVariable {
	var oValue domain.AdminPermissionVariable
	if oProto == nil {
		return oValue
	}

	oValue.Id = oProto.Id
	oValue.Key = oProto.Key
	oValue.Name = oProto.Name

	if oProto.Type != nil {
		iType := uint8(*oProto.Type)
		oValue.Type = &iType
	}

	return oValue
}
