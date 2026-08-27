package usecaseApplicationAnyAdminAuthentication

import (
	bootstrap "example/bootstrap"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	utility "example/internal/utility"
	pkgUtility "example/pkg/utility"
	"fmt"
)

type AuthenticatorUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	outputPortAnyModel.AdminUserModel
}

func NewAuthenticatorUsecase(oAminUserRepository outputPortAnyModel.AdminUserModel, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminAuthentication.AuthenticatorUsecase {
	return &AuthenticatorUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminUserModel:  oAminUserRepository,
	}
}

func (oSelf *AuthenticatorUsecase) SignIn(sName string, sPassword string, sSecret string) (string, error) {

	if sName == "" {
		return "", pkgUtility.NewDefaultError("name 不能為空", -1, 200)

	}

	if sPassword == "" {
		return "", pkgUtility.NewDefaultError("password 不能為空", -1, 200)
	}

	oAdminUser, err := oSelf.AdminUserModel.ShowOneByName(sName)

	if err != nil {
		return "", err
	}
	if oAdminUser == nil {
		return "", pkgUtility.NewDefaultError(sName+" 不存在", -2, 200)
	}

	sMd5 := utility.Md5(sPassword + bootstrap.CONFIG.TABLE.ADMIN_USER.PASSWORD)
	if oAdminUser.Password != sMd5 {
		return "", pkgUtility.NewDefaultError("密碼錯誤", -2, 200)
	}

	sAuthorization, err := oSelf.JwtHelper.Generate(int64(oAdminUser.Id), 0, map[string]any{}, sSecret)
	if err != nil {
		return "", pkgUtility.NewDefaultError("JWT 產生失敗", -2, 200)
	}

	return sAuthorization, nil
}

func (oSelf *AuthenticatorUsecase) Refresh(sJwt string, sSecret string) (string, error) {

	if sJwt == "" {
		return "", pkgUtility.NewDefaultError("JWT 不能為空", -1, 200)
	}

	fmt.Println("sJwt=", sJwt)

	oClaims, err := oSelf.JwtHelper.Parse(sJwt)
	if err != nil {
		return "", pkgUtility.NewDefaultError("JWT 無效", -2, 200)
	}

	iId := uint(oClaims.AdminUserId)
	oAdminUser, err := oSelf.AdminUserModel.ShowOneById(iId)
	fmt.Println("err=", err)

	if err != nil {
		return "", err
	}
	if oAdminUser == nil {
		return "", pkgUtility.NewDefaultError("AdminUser 不存在", -2, 200)
	}

	sAuthorization, err := oSelf.JwtHelper.Generate(oClaims.AdminUserId, oClaims.AppUserId, oClaims.Payload, sSecret)
	if err != nil {
		return "", pkgUtility.NewDefaultError("JWT 產生失敗", -2, 200)
	}

	return sAuthorization, nil
}
