package protoToDomain

import (
	domain "example/internal/domain"
	pb "example/pb"
)

func AdminPermissionGroupVariable(oProto *pb.AdminPermissionGroupVariable) domain.AdminPermissionGroupVariable {
	var oValue domain.AdminPermissionGroupVariable
	if oProto == nil {
		return oValue
	}

	oValue.Key = oProto.Key
	oValue.Name = oProto.Name

	return oValue
}
