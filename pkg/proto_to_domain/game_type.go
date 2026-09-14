package protoToDomain

import (
	domain "example/internal/domain"
	pb "example/pb"
)

func GameType(oProtoGameType *pb.GameType) domain.GameType {
	if oProtoGameType == nil {
		return domain.GameType{}
	}

	oGameType := domain.GameType{
		Id:        uint64(oProtoGameType.GetId()),
		ParentId:  uint64(oProtoGameType.GetParentId()),
		Key:       oProtoGameType.GetKey(),
		Name:      oProtoGameType.GetName(),
		CreatedAt: oProtoGameType.GetCreatedAt().AsTime(),
		UpdatedAt: oProtoGameType.GetUpdatedAt().AsTime(),
		DeletedAt: oProtoGameType.GetDeletedAt().AsTime(),
	}

	if oParent := oProtoGameType.GetParent(); oParent != nil {
		oParentDomain := GameType(oParent)
		oGameType.Parent = &oParentDomain
	}

	for _, oChild := range oProtoGameType.GetChildren() {
		oChildDomain := GameType(oChild)
		oGameType.Children = append(oGameType.Children, &oChildDomain)
	}

	return oGameType
}
