package authentication

import (
	"strings"

	bootstrap "example/bootstrap"
	inputApplicationTcp "example/internal/input/application/tcp"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	types "example/types"
)

type AuthenticatorHandler struct {
	*inputApplicationTcp.AbstractHandler
	AuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationTcp.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

// SignIn 簽名就是 pkg.TcpHandlerFunc，可以直接被 Tcp.HandleFunc 註冊。
// param 格式固定「name:password」。
func (oSelf *AuthenticatorHandler) SignIn(oReq types.TcpRequest) types.TcpResponse {
	aParts := strings.SplitN(oReq.Param, ":", 2)
	if len(aParts) != 2 {
		return types.TcpResponse{Code: -1, Message: "invalid param, expect \"name:password\""}
	}

	// NOTE: tcp carrier 目前沒有自己的 services.tcp.admin.jwt 設定，先借用 http 那組 secret。
	sAuthorization, err := oSelf.AuthenticatorUsecase.SignIn(aParts[0], aParts[1], bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.SECRET)
	if err != nil {
		return types.TcpResponse{Code: -1, Message: err.Error()}
	}

	return types.TcpResponse{Code: 1, Message: "成功登入", Result: sAuthorization}
}
