package domainToProto

import (
	domain "example/internal/domain"

	pbResource "example/pb/resource"
)

func GameValue(oGame *domain.GameVariable) *pbResource.GameValue {
	oValue := &pbResource.GameValue{
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
