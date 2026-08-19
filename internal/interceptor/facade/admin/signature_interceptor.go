package interceptor_facade_admin

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	bootstrap "example/bootstrap"
	utility "example/internal/utility"
)

type SignatureInterceptor struct {
	*AbstractInterceptor
}

func NewSignatureInterceptor(oInterceptor *AbstractInterceptor) *SignatureInterceptor {
	return &SignatureInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

// Handle 對齊 src/clients/facade.ts 的 oEncryptionInterceptor：client 端用
// [Ver, Version, K, Time, Salt] 以「|」串接後 MD5，放進 "signature" header，
// 這裡照同一個算法重算一次比對，簽名對不上就直接擋掉，不呼叫 fnHandler。
func (oSelf *SignatureInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(oContext context.Context, oRequest any, info *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (any, error) {

		oMd, _ := metadata.FromIncomingContext(oContext)

		fmt.Println(oMd, "md..")

		var sVer, sVersion, sK, sTime, sHeaderSignature string
		if aValues := oMd.Get("ver"); len(aValues) > 0 {
			sVer = aValues[0]
		}
		if aValues := oMd.Get("version"); len(aValues) > 0 {
			sVersion = aValues[0]
		}
		if aValues := oMd.Get("k"); len(aValues) > 0 {
			sK = aValues[0]
		}
		if aValues := oMd.Get("time"); len(aValues) > 0 {
			sTime = aValues[0]
		}
		if aValues := oMd.Get("signature"); len(aValues) > 0 {
			sHeaderSignature = aValues[0]
		}

		// NOTE: 不像 http 版本還要簽 s/o/p，facade 這邊要加密的欄位直接在 proto message
		// 裡原地換成密文，不會另外拆出 s/o/p 字串，所以只簽 Ver/Version/K/Time/Salt。
		aStrings := []string{sVer, sVersion, sK, sTime, bootstrap.CONFIG.SERVICES.FACADE.ADMIN.SALT}

		sStrings := strings.Join(aStrings, "|")

		fmt.Println(sStrings)

		sMd5Signature := utility.Md5(sStrings)

		if bootstrap.CONFIG.SERVICES.FACADE.ADMIN.SIGNATURE && sMd5Signature != sHeaderSignature {
			return nil, status.Error(codes.Unauthenticated, "簽名失敗")
		}

		return fnHandler(oContext, oRequest)
	}
}
