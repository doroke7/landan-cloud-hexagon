package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb/resource"
)

func AdminPermissionGroupVariable(oProto *pbResource.AdminPermissionGroupVariable) domain.AdminPermissionGroupVariable {
	var oValue domain.AdminPermissionGroupVariable
	if oProto == nil {
		return oValue
	}

	oValue.Key = oProto.Key
	oValue.Name = oProto.Name

	return oValue
}
