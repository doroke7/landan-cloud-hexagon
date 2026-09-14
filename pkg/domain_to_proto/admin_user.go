package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pbResource "example/pb/resource"
)

func AdminUser(oDomain *domain.AdminUser) *pbResource.AdminUser {
	if oDomain == nil {
		return nil
	}

	aAdminRoles := make([]*pbResource.AdminRole, 0, len(oDomain.AdminRoles))
	for i := range oDomain.AdminRoles {
		aAdminRoles = append(aAdminRoles, AdminRole(&oDomain.AdminRoles[i]))
	}

	return &pbResource.AdminUser{
		Id:         oDomain.Id,
		Name:       oDomain.Name,
		Password:   oDomain.Password,
		CreatedAt:  timestamppb.New(oDomain.CreatedAt),
		UpdatedAt:  timestamppb.New(oDomain.UpdatedAt),
		DeletedAt:  timestamppb.New(oDomain.DeletedAt),
		AdminRoles: aAdminRoles,
	}
}
