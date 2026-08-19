package authentication

import (
	"encoding/json"

	bootstrap "example/bootstrap"
	inputApplicationWebsocket "example/internal/input/application/websocket"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	types "example/types"
)

type signInParam struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthenticatorHandler struct {
	*inputApplicationWebsocket.AbstractHandler
	AuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationWebsocket.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) SignIn(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse {
	var oParam signInParam
	if err := json.Unmarshal(oReq.Param, &oParam); err != nil {
		return types.WebsocketResponse{Type: "normal", Code: -1, Message: "invalid param, expect {\"name\":..., \"password\":...}"}
	}

	sAuthorization, err := oSelf.AuthenticatorUsecase.SignIn(oParam.Name, oParam.Password, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.JWT.SECRET)
	if err != nil {
		return types.WebsocketResponse{Type: "normal", Code: -1, Message: err.Error()}
	}

	if oConn != nil {
		_ = oConn.Push("Admin.Authentication.Authenticator.SignIn.Welcome", map[string]any{"name": oParam.Name})
	}

	// sAuthorization 是 string，json.Marshal 一個 string 不會出錯，可以放心略過 err。
	aResult, _ := json.Marshal(sAuthorization)
	return types.WebsocketResponse{Type: "ack", Code: 1, Message: "成功登入", Result: aResult}
}
