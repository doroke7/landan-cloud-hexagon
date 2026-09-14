package protoToDomain

import (
	domain "example/internal/domain"
	pbResource "example/pb"
)

func GameVariable(oVariable *pbResource.GameVariable) domain.GameVariable {
	var oValue domain.GameVariable
	if oVariable == nil {
		return oValue
	}

	oValue.Key = oVariable.Key
	oValue.Name = oVariable.Name
	oValue.Description = oVariable.Description

	if oVariable.GameTypeId != nil {
		iGameTypeId := uint64(*oVariable.GameTypeId)
		oValue.GameTypeId = &iGameTypeId
	}

	return oValue
}
