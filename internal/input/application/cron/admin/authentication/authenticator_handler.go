package authentication

import (
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	pkg "example/pkg"

	inputApplicationCron "example/internal/input/application/cron"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
)

type AuthenticatorHandler struct {
	*inputApplicationCron.AbstractHandler
	AuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationCron.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) SignIn() {
	// NOTE: cron carrier 目前沒有自己的 services.cron.admin.jwt 設定，先借用 http 那組 secret。
	sAuthorization, err := oSelf.AuthenticatorUsecase.SignIn("tom", "secret", bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.SECRET)
	if err != nil {
		pkg.Logger(pkg.Cron).Error("SignIn 失敗",
			zap.Error(err),
		)
		return
	}

	pkg.Logger(pkg.Cron).Info("SignIn 成功",
		zap.String("authorization", sAuthorization),
	)
}
