package usecase

import (
	helper "example/internal/helper"
)

type AbstractUsecase struct {
	*helper.AesHelper
	*helper.JwtHelper
}

func NewAbstractUsecase(oAesHelper *helper.AesHelper, oJwtHelper *helper.JwtHelper) *AbstractUsecase {
	return &AbstractUsecase{
		AesHelper: oAesHelper,
		JwtHelper: oJwtHelper,
	}
}
