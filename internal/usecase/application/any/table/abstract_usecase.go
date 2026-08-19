package usecase

import (
	helper "example/internal/helper"
)

type AbstractUsecase struct {
	*helper.AesHelper
}

func NewAbstractUsecase(oAesHelper *helper.AesHelper) *AbstractUsecase {
	return &AbstractUsecase{
		AesHelper: oAesHelper,
	}
}
