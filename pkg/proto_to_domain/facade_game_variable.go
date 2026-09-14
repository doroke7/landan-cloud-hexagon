package protoToDomain

import (
	domain "example/internal/domain"
	pbFacade "example/pb/facade"
)

func FacadeGameVariable(oVariable *pbFacade.GameVariable) *domain.GameVariable {
	oDomain := &domain.GameVariable{}
	if oVariable == nil {
		return oDomain
	}

	if oVariable.GameTypeId != nil {
		iGameTypeIdValue := oVariable.GetGameTypeId()
		iGameTypeId := uint64(iGameTypeIdValue)
		oDomain.GameTypeId = &iGameTypeId
	}

	oDomain.Key = oVariable.Key
	oDomain.Name = oVariable.Name
	oDomain.Description = oVariable.Description

	return oDomain
}
