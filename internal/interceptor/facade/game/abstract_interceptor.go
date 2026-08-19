package interceptor_facade_game

import (
	"context"

	"google.golang.org/grpc"

	helper "example/internal/helper"
)

type AbstractInterceptor struct {
	JwtHelper *helper.JwtHelper
	AesHelper *helper.AesHelper
	RsaHelper *helper.RsaHelper
}

func NewAbstractInterceptor(oJwtHelper *helper.JwtHelper, oAesHelper *helper.AesHelper, oRsaHelper *helper.RsaHelper) *AbstractInterceptor {
	return &AbstractInterceptor{
		JwtHelper: oJwtHelper,
		AesHelper: oAesHelper,
		RsaHelper: oRsaHelper,
	}
}

func (oSelf *AbstractInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(oContex context.Context, oRequest any, info *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (any, error) {
		return fnHandler(oContex, oRequest)
	}
}
