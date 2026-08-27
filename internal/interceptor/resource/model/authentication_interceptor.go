package interceptorResourceModel

import (
	"context"
	"strings"

	bootstrap "example/bootstrap"
	pkgUtility "example/pkg/utility"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthenticationInterceptor struct {
	*AbstractInterceptor
}

func NewAuthenticationInterceptor(oInterceptor *AbstractInterceptor) *AuthenticationInterceptor {
	return &AuthenticationInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

// Handle 驗證 facade -> resource 呼叫的 Basic Auth
func (oSelf *AuthenticationInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(oCtx context.Context, oReq any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		oMetadata, bMetaOb := metadata.FromIncomingContext(oCtx)

		if bMetaOb {
			aAuthotizations := oMetadata.Get("authorization")
			sAuthotizations := strings.Join(aAuthotizations, "")

			sUser := bootstrap.CONFIG.SERVICES.RESOURCE.USER
			sPassword := bootstrap.CONFIG.SERVICES.RESOURCE.PASSWORD

			sAuthotization := "Basic " + pkgUtility.Base64Encode(sUser+":"+sPassword)

			if sAuthotizations != sAuthotization {
				return nil, status.Error(codes.PermissionDenied, "resource 密碼錯誤")
			}
		}

		return handler(oCtx, oReq)
	}
}
