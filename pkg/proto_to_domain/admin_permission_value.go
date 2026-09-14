package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb/resource"
)

func AdminPermissionValue(oProto *pbResource.AdminPermissionValue) *domain.AdminPermissionValue {
	if oProto == nil {
		return nil
	}

	var oType *uint8
	if oProto.Type != nil {
		iType := uint8(*oProto.Type)
		oType = &iType
	}

	oDomain := &domain.AdminPermissionValue{
		Id:   oProto.Id,
		Type: oType,
		Key:  oProto.Key,
		Name: oProto.Name,
	}

	return oDomain
}
