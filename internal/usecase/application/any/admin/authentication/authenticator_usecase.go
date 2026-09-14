package usecaseApplicationAnyAdminAuthentication

import (
	bootstrap "example/bootstrap"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	pkgUtility "example/pkg/utility"
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
		oNameEmptyError := pkgUtility.NewDefaultError("name must not be empty", -1, 200)

		return "", oNameEmptyError
	}

	if sPassword == "" {
		oPasswordEmptyError := pkgUtility.NewDefaultError("password must not be empty", -1, 200)

		return "", oPasswordEmptyError
	}

	oAdminUser, err := oSelf.AdminUserModel.ShowOneByName(sName)

	if err != nil {
		return "", err
	}
	if oAdminUser == nil {
		sNotExistMessage := sName + " does not exist"
		oNotExistError := pkgUtility.NewDefaultError(sNotExistMessage, -2, 200)

		return "", oNotExistError
	}

	sMd5 := pkgUtility.Md5(sPassword + bootstrap.CONFIG.TABLE.ADMIN_USER.PASSWORD)
	if oAdminUser.Password != sMd5 {
		oIncorrectPasswordError := pkgUtility.NewDefaultError("incorrect password", -2, 200)

		return "", oIncorrectPasswordError
	}

	sAuthorization, err := oSelf.JwtHelper.Generate(int64(oAdminUser.Id), 0, map[string]any{}, sSecret)
	if err != nil {
		oGenerateError := pkgUtility.NewDefaultError("JWT generation failed", -2, 200)

		return "", oGenerateError
	}

	return sAuthorization, nil
}

func (oSelf *AuthenticatorUsecase) Refresh(sJwt string, sSecret string) (string, error) {

	if sJwt == "" {
		oJwtEmptyError := pkgUtility.NewDefaultError("JWT must not be empty", -1, 200)

		return "", oJwtEmptyError
	}

	oClaims, err := oSelf.JwtHelper.Parse(sJwt)
	if err != nil {
		oInvalidJwtError := pkgUtility.NewDefaultError("invalid JWT", -2, 200)

		return "", oInvalidJwtError
	}

	iId := uint64(oClaims.AdminUserId)
	oAdminUser, err := oSelf.AdminUserModel.ShowOneById(iId)

	if err != nil {
		return "", err
	}
	if oAdminUser == nil {
		oNotExistError := pkgUtility.NewDefaultError("AdminUser does not exist", -2, 200)

		return "", oNotExistError
	}

	sAuthorization, err := oSelf.JwtHelper.Generate(oClaims.AdminUserId, oClaims.AppUserId, oClaims.Payload, sSecret)
	if err != nil {
		oGenerateError := pkgUtility.NewDefaultError("JWT generation failed", -2, 200)

		return "", oGenerateError
	}

	return sAuthorization, nil
}
