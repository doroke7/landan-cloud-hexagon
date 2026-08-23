package interceptorFacadeAdmin

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	bootstrap "example/bootstrap"
)

type EncryptionInterceptor struct {
	*AbstractInterceptor
}

func NewEncryptionInterceptor(oInterceptor *AbstractInterceptor) *EncryptionInterceptor {
	return &EncryptionInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

// Handle 呼叫 handler 前把 *authHolder 塞進 context；handler 執行完之後（after 階段）
// 讀出 SignIn 寫進 holder 的明文 authorization，用 AES 加密，寫進 "A" header。
func (oSelf *EncryptionInterceptor) Handle() grpc.UnaryServerInterceptor {

	// gRPC 的 context 是值
	// go-Gin 的 context 是指標
	return func(oContext context.Context, oRequest any, oServerInfo *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (oResponse any, oErr error) {

		oAuthorization := new(string)
		oContext = context.WithValue(oContext, "a", oAuthorization)

		// 由於 fnHandler func 本身是傳值，也就是值拷貝， 相當隱含 傳入了一份 子 context進去
		// 由於 context 是不可變，也就是 子 context 改變不會改變父 context
		oResponse, oErr = fnHandler(oContext, oRequest)
		sAuthorization := *oAuthorization
		if oErr != nil {
			return oResponse, oErr
		}

		if sAuthorization == "" {
			return oResponse, oErr
		}

		sA, oEncryptErr := oSelf.AesHelper.Encrypt(
			sAuthorization,
			bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.KEY,
			bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.IV,
		)
		if oEncryptErr != nil {
			return nil, oEncryptErr
		}

		if oHeaderErr := grpc.SetHeader(oContext, metadata.Pairs("A", sA)); oHeaderErr != nil {
			return nil, oHeaderErr
		}

		return oResponse, oErr
	}

}
