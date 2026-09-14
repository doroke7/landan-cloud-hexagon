package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pbResource "example/pb/resource"
)

func AdminRole(oDomain *domain.AdminRole) *pbResource.AdminRole {
	if oDomain == nil {
		return nil
	}

	return &pbResource.AdminRole{
		Id:        oDomain.Id,
		Key:       oDomain.Key,
		Name:      oDomain.Name,
		CreatedAt: timestamppb.New(oDomain.CreatedAt),
		UpdatedAt: timestamppb.New(oDomain.UpdatedAt),
		DeletedAt: timestamppb.New(oDomain.DeletedAt),
	}
}
