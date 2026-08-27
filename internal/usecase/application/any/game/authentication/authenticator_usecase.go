package usecaseApplicationAnyGameAuthentication

import (
	bootstrap "example/bootstrap"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyGame "example/internal/usecase/application/any/game"
	usecasePortAnyGameAuthentication "example/internal/usecase/port/any/game/authentication"
	utility "example/internal/utility"
	pkgUtility "example/pkg/utility"
	"fmt"
)

type AuthenticatorUsecase struct {
	*usecaseApplicationAnyGame.AbstractUsecase
	outputPortAnyModel.AppUserModel
}

func NewAuthenticatorUsecase(oAppUserRepository outputPortAnyModel.AppUserModel, oAbstractUsecase *usecaseApplicationAnyGame.AbstractUsecase) usecasePortAnyGameAuthentication.AuthenticatorUsecase {
	return &AuthenticatorUsecase{
		AbstractUsecase: oAbstractUsecase,
		AppUserModel:    oAppUserRepository,
	}
}

func (oSelf *AuthenticatorUsecase) LogIn(sName string, sPassword string, sSecret string) (string, error) {

	if sName == "" {
		return "", pkgUtility.NewDefaultError("name 不能為空", -1, 200)
	}

	if sPassword == "" {
		return "", pkgUtility.NewDefaultError("password 不能為空", -1, 200)
	}

	oAppUser, oErr := oSelf.AppUserModel.ShowOneByName(sName)
	if oErr != nil {
		return "", oErr
	}
	if oAppUser == nil {
		return "", pkgUtility.NewDefaultError(sName+" 不存在", -2, 200)
	}

	sMd5 := utility.Md5(sPassword + bootstrap.CONFIG.TABLE.APP_USER.PASSWORD)
	fmt.Println("sMd5=", sMd5)
	if oAppUser.Password != sMd5 {
		return "", pkgUtility.NewDefaultError("密碼錯誤", -2, 200)
	}

	sAuthorization, oErr := oSelf.JwtHelper.Generate(0, int64(oAppUser.Id), map[string]any{}, sSecret)
	if oErr != nil {
		return "", pkgUtility.NewDefaultError("JWT 產生失敗", -2, 200)
	}

	return sAuthorization, nil
}

func (oSelf *AuthenticatorUsecase) Refresh(sJwt string, sSecret string) (string, error) {

	if sJwt == "" {
		return "", pkgUtility.NewDefaultError("JWT 不能為空", -1, 200)
	}

	oClaims, oErr := oSelf.JwtHelper.Parse(sJwt)
	if oErr != nil {
		return "", pkgUtility.NewDefaultError("JWT 無效", -2, 200)
	}

	iId := uint(oClaims.AppUserId)
	oAppUser, oErr := oSelf.AppUserModel.ShowOneById(iId)
	if oErr != nil {
		return "", oErr
	}
	if oAppUser == nil {
		return "", pkgUtility.NewDefaultError("AppUser 不存在", -2, 200)
	}

	sAuthorization, oErr := oSelf.JwtHelper.Generate(oClaims.AdminUserId, oClaims.AppUserId, oClaims.Payload, sSecret)
	if oErr != nil {
		return "", pkgUtility.NewDefaultError("JWT 產生失敗", -2, 200)
	}

	return sAuthorization, nil
}
