package interceptorFacadeAdmin

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	bootstrap "example/bootstrap"
	pkgUtility "example/pkg/utility"
)

type DecryptionInterceptor struct {
	*AbstractInterceptor
}

func NewDecryptionInterceptor(oInterceptor *AbstractInterceptor) *DecryptionInterceptor {
	return &DecryptionInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

// encryptedFieldOptionNumber 對應 proto 裡 `extend google.protobuf.FieldOptions { bool encrypted = 50001; }`
// 這個 extension number，跟 src/clients/facade.ts 的 getOption(oField, encrypted) 是同一個欄位選項。
// 用 number 判斷、不綁死某一個 proto 檔案產生的 Go extension 符號（例如 pb_facade_admin_authentication.E_Encrypted），
// 這樣 admin 底下任何訊息只要用了同一個 extension number，都能被這個攔截器認得。
const encryptedFieldOptionNumber = 50001

// Handle 對齊 src/clients/facade.ts 的 oEncryptionInterceptor：client 端把
// proto message 裡標記 (encrypted)=true 的欄位，用隨機產生的 AES key/iv 加密，
// key/iv 打包成 JSON 後用 RSA 公鑰加密放進 "k" header。這裡反過來做：
// 先用 RSA 私鑰解出 key/iv，再把訊息裡每個標記加密的欄位解密回明文，
// 原地覆寫回 proto message，後面的 handler 拿到的就是明文。
func (oSelf *DecryptionInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(oContext context.Context, oRequest any, info *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (any, error) {

		oMessage, bOk := oRequest.(proto.Message)
		if !bOk {
			return fnHandler(oContext, oRequest)
		}

		oReflect := oMessage.ProtoReflect()
		aEncryptedFields := findEncryptedFields(oReflect.Descriptor())
		if len(aEncryptedFields) == 0 {
			return fnHandler(oContext, oRequest)
		}

		oMd, _ := metadata.FromIncomingContext(oContext)

		var sK string
		if aValues := oMd.Get("k"); len(aValues) > 0 {
			sK = aValues[0]
		}
		if sK == "" {
			return nil, status.Error(codes.InvalidArgument, "missing k header")
		}

		sKeys, oErr := oSelf.RsaHelper.Decrypt(sK, bootstrap.CONFIG.SERVICES.FACADE.ADMIN.PRIVATE_KEY)
		if oErr != nil {
			return nil, status.Error(codes.InvalidArgument, "key decryption failed")
		}

		oKeys, oErr := pkgUtility.JsonDecode[struct {
			Key string `json:"key"`
			Iv  string `json:"iv"`
		}](sKeys)
		if oErr != nil {
			return nil, status.Error(codes.InvalidArgument, "key format error")
		}

		for _, oField := range aEncryptedFields {
			sValue := oReflect.Get(oField).String()
			if sValue == "" {
				continue
			}

			// NOTE: src/clients/facade.ts 呼叫 aesHelper.encrypt(..., { format: 'base64' }) 時，
			// 本來想把標準 base64 的 +/ 換成 url-safe 的 -_，但那行 .replace(...) 沒把結果指定回
			// 變數，是個沒生效的 no-op，實際送出來的還是標準 base64（帶 +/）。這裡統一先轉成
			// url-safe，AesHelper.Decrypt 固定用 RawURLEncoding 解，兩種格式都吃得下。
			sNormalized := strings.NewReplacer("+", "-", "/", "_").Replace(sValue)

			fmt.Println("oKeys 90", oKeys)
			fmt.Println("sNormalized ", sNormalized)

			sDecrypted, oErr := oSelf.AesHelper.Decrypt(sNormalized, oKeys.Key, oKeys.Iv)
			if oErr != nil {
				return nil, status.Error(codes.InvalidArgument, "content decryption failed")
			}

			oReflect.Set(oField, protoreflect.ValueOfString(sDecrypted))
		}

		return fnHandler(oContext, oRequest)
	}
}

func findEncryptedFields(oDescriptor protoreflect.MessageDescriptor) []protoreflect.FieldDescriptor {
	aFields := oDescriptor.Fields()
	aResult := make([]protoreflect.FieldDescriptor, 0, aFields.Len())

	for i := 0; i < aFields.Len(); i++ {
		oField := aFields.Get(i)
		if oField.Kind() != protoreflect.StringKind {
			continue
		}
		if isEncryptedField(oField) {
			aResult = append(aResult, oField)
		}
	}

	return aResult
}

func isEncryptedField(oField protoreflect.FieldDescriptor) bool {
	oOptions, bOk := oField.Options().(*descriptorpb.FieldOptions)
	if !bOk || oOptions == nil {
		return false
	}

	bEncrypted := false
	proto.RangeExtensions(oOptions, func(oExtType protoreflect.ExtensionType, mValue any) bool {
		if oExtType.TypeDescriptor().Number() != encryptedFieldOptionNumber {
			return true
		}
		if bValue, bOk := mValue.(bool); bOk {
			bEncrypted = bValue
		}
		return true
	})

	return bEncrypted
}
