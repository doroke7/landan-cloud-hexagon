package domainToProto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"

	pb "example/pb"
)

func AdminRole(oDomain *domain.AdminRole) *pb.AdminRole {
	if oDomain == nil {
		return nil
	}

	return &pb.AdminRole{
		Id:        oDomain.Id,
		Key:       oDomain.Key,
		Name:      oDomain.Name,
		CreatedAt: timestamppb.New(oDomain.CreatedAt),
		UpdatedAt: timestamppb.New(oDomain.UpdatedAt),
		DeletedAt: timestamppb.New(oDomain.DeletedAt),
	}
}
