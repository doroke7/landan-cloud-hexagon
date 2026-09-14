package protoToDomain

import (
	domain "example/internal/domain"
	pb "example/pb"
)

func AdminUser(oProto *pb.AdminUser) domain.AdminUser {
	if oProto == nil {
		return domain.AdminUser{}
	}

	aAdminRoles := make([]*domain.AdminRole, 0, len(oProto.GetAdminRoles()))
	for _, oProtoAdminRole := range oProto.GetAdminRoles() {
		oDomainAdminRole := AdminRole(oProtoAdminRole)
		aAdminRoles = append(aAdminRoles, &oDomainAdminRole)
	}

	return domain.AdminUser{
		Id:         uint64(oProto.GetId()),
		Name:       oProto.GetName(),
		Password:   oProto.GetPassword(),
		CreatedAt:  oProto.GetCreatedAt().AsTime(),
		UpdatedAt:  oProto.GetUpdatedAt().AsTime(),
		DeletedAt:  oProto.GetDeletedAt().AsTime(),
		AdminRoles: aAdminRoles,
	}
}
