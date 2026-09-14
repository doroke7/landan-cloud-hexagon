package domainToProto

import (
	domain "example/internal/domain"
	pb "example/pb"
)

func AdminPermissionVariable(oDomainAdminPermissionValue *domain.AdminPermissionVariable) *pb.AdminPermissionVariable {
	if oDomainAdminPermissionValue == nil {
		return nil
	}

	oProtoAdminPermissionValue := &pb.AdminPermissionVariable{
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
