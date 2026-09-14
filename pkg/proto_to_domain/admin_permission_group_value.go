package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb/resource"
)

func AdminPermissionGroupValue(oProto *pbResource.AdminPermissionGroupValue) domain.AdminPermissionGroupVariable {
	var oValue domain.AdminPermissionGroupVariable
	if oProto == nil {
		return oValue
	}

	oValue.Key = oProto.Key
	oValue.Name = oProto.Name

	return oValue
}
