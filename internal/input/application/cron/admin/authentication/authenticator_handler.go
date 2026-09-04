package inputApplicationCronAdminAuthentication

import (
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	pkgUtility "example/pkg/utility"

	inputApplicationCron "example/internal/input/application/cron"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
)

type AuthenticatorHandler struct {
	*inputApplicationCron.AbstractHandler
	AdminAuthenticationAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationCron.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:                         oAbstractHandler,
		AdminAuthenticationAuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) SignIn() {
	// NOTE: cron carrier 目前沒有自己的 services.cron.admin.jwt 設定，先借用 http 那組 secret。
	sAuthorization, err := oSelf.AdminAuthenticationAuthenticatorUsecase.SignIn("tom", "secret", bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.SECRET)
	if err != nil {
		pkgUtility.Logger(pkgUtility.Cron).Error("SignIn 失敗",
			zap.Error(err),
		)
		return
	}

	pkgUtility.Logger(pkgUtility.Cron).Info("SignIn 成功",
		zap.String("authorization", sAuthorization),
	)
}
