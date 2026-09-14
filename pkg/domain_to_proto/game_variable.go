package domainToProto

import (
	domain "example/internal/domain"

	pbResource "example/pb"
)

func GameVariable(oGame *domain.GameVariable) *pbResource.GameVariable {
	oValue := &pbResource.GameVariable{
		Key:         oGame.Key,
		Name:        oGame.Name,
		Description: oGame.Description,
	}

	if oGame.GameTypeId != nil {
		iGameTypeId := uint64(*oGame.GameTypeId)
		oValue.GameTypeId = &iGameTypeId
	}

	return oValue
}
