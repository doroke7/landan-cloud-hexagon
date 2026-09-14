package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	bootstrap "example/bootstrap"
)

type JwtHelper struct {
	*AbstractHelper
}

func NewJwtHelper(oAbstractHelper *AbstractHelper) *JwtHelper {
	return &JwtHelper{
		AbstractHelper: oAbstractHelper,
	}
}

// JwtClaims 自定義 claims，Payload 可存任意業務資料
type JwtClaims struct {
	jwt.RegisteredClaims
	AdminUserId int64          `json:"admin_user_id"`
	AppUserId   int64          `json:"app_user_id"`
	Payload     map[string]any `json:"payload"`
}

// Generate 簽發 JWT 後用 AES 加密，回傳加密後的 token。
// sSecret 由呼叫端決定要用哪個 carrier 的 SERVICES.<CARRIER>.ADMIN.JWT.SECRET，
// JwtHelper 本身不綁定任何一個 carrier。
func (oSelf *JwtHelper) Generate(nAdminUserId int64, nAppUserId int64, oPayload map[string]any, sSecret string) (string, error) {
	oNow := time.Now()
	oExpiresAtTime := oNow.Add(24 * time.Hour)
	oIssuedAt := jwt.NewNumericDate(oNow)
	oExpiresAt := jwt.NewNumericDate(oExpiresAtTime)

	oClaims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  oIssuedAt,
			ExpiresAt: oExpiresAt,
		},
		AdminUserId: nAdminUserId,
		AppUserId:   nAppUserId,
		Payload:     oPayload,
	}

	oToken := jwt.NewWithClaims(jwt.SigningMethodHS256, oClaims)
	sJwt, oErr := oToken.SignedString([]byte(sSecret))
	if oErr != nil {
		return "", oErr
	}

	return sJwt, nil
}

// Parse 先 AES 解密，再驗證並解析 JWT，回傳 JwtClaims
func (oSelf *JwtHelper) Parse(sJwt string) (*JwtClaims, error) {

	oToken, oErr := jwt.ParseWithClaims(sJwt, &JwtClaims{}, func(oT *jwt.Token) (any, error) {
		if _, bOk := oT.Method.(*jwt.SigningMethodHMAC); !bOk {
			oErr := errors.New("unexpected signing method")

			return nil, oErr
		}
		return []byte(bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.SECRET), nil
	})
	if oErr != nil {
		return nil, oErr
	}

	oClaims, bOk := oToken.Claims.(*JwtClaims)
	if !bOk || !oToken.Valid {
		oErr := errors.New("invalid token")

		return nil, oErr
	}

	return oClaims, nil
}

// Refresh 驗證舊 token 後，以相同 Payload 簽發新 token
func (oSelf *JwtHelper) Refresh(sToken string, sSecret string) (string, error) {
	oClaims, oErr := oSelf.Parse(sToken)
	if oErr != nil {
		return "", oErr
	}

	sNewToken, oGenerateErr := oSelf.Generate(oClaims.AdminUserId, oClaims.AppUserId, oClaims.Payload, sSecret)

	return sNewToken, oGenerateErr
}
