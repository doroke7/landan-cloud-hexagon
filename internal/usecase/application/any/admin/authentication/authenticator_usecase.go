package usecaseApplicationAnyAdminAuthentication

import (
	bootstrap "example/bootstrap"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
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
		return "", pkgUtility.NewDefaultError("name must not be empty", -1, 200)

	}

	if sPassword == "" {
		return "", pkgUtility.NewDefaultError("password must not be empty", -1, 200)
	}

	oAdminUser, err := oSelf.AdminUserModel.ShowOneByName(sName)

	if err != nil {
		return "", err
	}
	if oAdminUser == nil {
		return "", pkgUtility.NewDefaultError(sName+" does not exist", -2, 200)
	}

	sMd5 := pkgUtility.Md5(sPassword + bootstrap.CONFIG.TABLE.ADMIN_USER.PASSWORD)
	if oAdminUser.Password != sMd5 {
		return "", pkgUtility.NewDefaultError("incorrect password", -2, 200)
	}

	sAuthorization, err := oSelf.JwtHelper.Generate(int64(oAdminUser.Id), 0, map[string]any{}, sSecret)
	if err != nil {
		return "", pkgUtility.NewDefaultError("JWT generation failed", -2, 200)
	}

	return sAuthorization, nil
}

func (oSelf *AuthenticatorUsecase) Refresh(sJwt string, sSecret string) (string, error) {

	if sJwt == "" {
		return "", pkgUtility.NewDefaultError("JWT must not be empty", -1, 200)
	}

	fmt.Println("sJwt=", sJwt)

	oClaims, err := oSelf.JwtHelper.Parse(sJwt)
	if err != nil {
		return "", pkgUtility.NewDefaultError("invalid JWT", -2, 200)
	}

	iId := uint(oClaims.AdminUserId)
	oAdminUser, err := oSelf.AdminUserModel.ShowOneById(iId)
	fmt.Println("err=", err)

	if err != nil {
		return "", err
	}
	if oAdminUser == nil {
		return "", pkgUtility.NewDefaultError("AdminUser does not exist", -2, 200)
	}

	sAuthorization, err := oSelf.JwtHelper.Generate(oClaims.AdminUserId, oClaims.AppUserId, oClaims.Payload, sSecret)
	if err != nil {
		return "", pkgUtility.NewDefaultError("JWT generation failed", -2, 200)
	}

	return sAuthorization, nil
}
