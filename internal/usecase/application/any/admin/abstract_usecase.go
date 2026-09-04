package usecaseApplicationAnyAdmin

import (
	helper "example/internal/helper"
)

type AbstractUsecase struct {
	*helper.AesHelper
	*helper.JwtHelper
	*helper.ValidatorHelper
}

func NewAbstractUsecase(oAesHelper *helper.AesHelper, oJwtHelper *helper.JwtHelper, oValidatorHelper *helper.ValidatorHelper) *AbstractUsecase {
	return &AbstractUsecase{
		AesHelper:       oAesHelper,
		JwtHelper:       oJwtHelper,
		ValidatorHelper: oValidatorHelper,
	}
}
