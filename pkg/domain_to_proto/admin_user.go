package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pb "example/pb"
)

func AdminUser(oDomain *domain.AdminUser) *pb.AdminUser {
	if oDomain == nil {
		return nil
	}

	aAdminRoles := make([]*pb.AdminRole, 0, len(oDomain.AdminRoles))
	for i := range oDomain.AdminRoles {
		aAdminRoles = append(aAdminRoles, AdminRole(&oDomain.AdminRoles[i]))
	}

	return &pb.AdminUser{
		Id:         oDomain.Id,
		Name:       oDomain.Name,
		Password:   oDomain.Password,
		CreatedAt:  timestamppb.New(oDomain.CreatedAt),
		UpdatedAt:  timestamppb.New(oDomain.UpdatedAt),
		DeletedAt:  timestamppb.New(oDomain.DeletedAt),
		AdminRoles: aAdminRoles,
	}
}
