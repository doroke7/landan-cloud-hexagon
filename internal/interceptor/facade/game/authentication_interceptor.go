package interceptor_facade_game

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const AdminUserIDKey contextKey = "admin_user_id"

type AuthenticationInterceptor struct {
	*AbstractInterceptor
}

func NewAuthenticationInterceptor(oInterceptor *AbstractInterceptor) *AuthenticationInterceptor {
	return &AuthenticationInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

func (oSelf *AuthenticationInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(oContext context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		oMd, bOk := metadata.FromIncomingContext(oContext)
		if !bOk {
			return nil, status.Error(codes.Unauthenticated, "缺少認證資訊")
		}

		aValues := oMd.Get("authorization")
		if len(aValues) == 0 {
			return nil, status.Error(codes.Unauthenticated, "缺少 authorization header")
		}

		sToken := strings.TrimPrefix(aValues[0], "Bearer ")

		oClaims, oErr := oSelf.JwtHelper.Parse(sToken)
		if oErr != nil {
			return nil, status.Error(codes.Unauthenticated, "token 無效或已過期")
		}

		oContext = context.WithValue(oContext, AdminUserIDKey, oClaims.AdminUserId)
		return handler(oContext, req)
	}
}
