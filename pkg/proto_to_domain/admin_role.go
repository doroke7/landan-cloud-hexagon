package protoToDomain

import (
	domain "example/internal/domain"
	pb "example/pb"
)

func AdminRole(oProto *pb.AdminRole) domain.AdminRole {
	if oProto == nil {
		return domain.AdminRole{}
	}

	oAdminRole := domain.AdminRole{
		Id:        uint64(oProto.GetId()),
		Key:       oProto.GetKey(),
		Name:      oProto.GetName(),
		CreatedAt: oProto.GetCreatedAt().AsTime(),
		UpdatedAt: oProto.GetUpdatedAt().AsTime(),
		DeletedAt: oProto.GetDeletedAt().AsTime(),
	}

	return oAdminRole
}
