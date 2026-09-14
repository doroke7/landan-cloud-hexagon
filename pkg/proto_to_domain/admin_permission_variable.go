package protoToDomain

import (
	domain "example/internal/domain"
	pb "example/pb"
)

func AdminPermissionVariable(oProto *pb.AdminPermissionVariable) domain.AdminPermissionVariable {
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
