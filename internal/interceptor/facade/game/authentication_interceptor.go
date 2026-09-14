package interceptorFacadeGame

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pkgUtility "example/pkg/utility"
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
			oMissingAuthenticationInfoError := pkgUtility.NewDefaultError("missing authentication info", -2, 200)

			return nil, oMissingAuthenticationInfoError
		}

		aValues := oMd.Get("authorization")
		if len(aValues) == 0 {
			oMissingAuthorizationHeaderError := pkgUtility.NewDefaultError("missing authorization header", -2, 200)

			return nil, oMissingAuthorizationHeaderError
		}

		sToken := strings.TrimPrefix(aValues[0], "Bearer ")

		oClaims, oErr := oSelf.JwtHelper.Parse(sToken)
		if oErr != nil {
			oInvalidTokenError := pkgUtility.NewDefaultError("invalid or expired token", -2, 200)

			return nil, oInvalidTokenError
		}

		oContext = context.WithValue(oContext, AdminUserIDKey, oClaims.AdminUserId)
		return handler(oContext, req)
	}
}
