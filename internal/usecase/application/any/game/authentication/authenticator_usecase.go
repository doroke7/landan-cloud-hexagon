package usecaseApplicationAnyGameAuthentication

import (
	bootstrap "example/bootstrap"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyGame "example/internal/usecase/application/any/game"
	usecasePortAnyGameAuthentication "example/internal/usecase/port/any/game/authentication"
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
		oNameEmptyError := pkgUtility.NewDefaultError("name must not be empty", -1, 200)

		return "", oNameEmptyError
	}

	if sPassword == "" {
		oPasswordEmptyError := pkgUtility.NewDefaultError("password must not be empty", -1, 200)

		return "", oPasswordEmptyError
	}

	oAppUser, oErr := oSelf.AppUserModel.ShowOneByName(sName)
	if oErr != nil {
		return "", oErr
	}
	if oAppUser == nil {
		sNotExistMessage := sName + " does not exist"
		oNotExistError := pkgUtility.NewDefaultError(sNotExistMessage, -2, 200)

		return "", oNotExistError
	}

	sMd5 := pkgUtility.Md5(sPassword + bootstrap.CONFIG.TABLE.APP_USER.PASSWORD)
	fmt.Println("sMd5=", sMd5)
	if oAppUser.Password != sMd5 {
		oIncorrectPasswordError := pkgUtility.NewDefaultError("incorrect password", -2, 200)

		return "", oIncorrectPasswordError
	}

	sAuthorization, oErr := oSelf.JwtHelper.Generate(0, int64(oAppUser.Id), map[string]any{}, sSecret)
	if oErr != nil {
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

	oClaims, oErr := oSelf.JwtHelper.Parse(sJwt)
	if oErr != nil {
		oInvalidJwtError := pkgUtility.NewDefaultError("invalid JWT", -2, 200)

		return "", oInvalidJwtError
	}

	iId := uint(oClaims.AppUserId)
	oAppUser, oErr := oSelf.AppUserModel.ShowOneById(iId)
	if oErr != nil {
		return "", oErr
	}
	if oAppUser == nil {
		oNotExistError := pkgUtility.NewDefaultError("AppUser does not exist", -2, 200)

		return "", oNotExistError
	}

	sAuthorization, oErr := oSelf.JwtHelper.Generate(oClaims.AdminUserId, oClaims.AppUserId, oClaims.Payload, sSecret)
	if oErr != nil {
		oGenerateError := pkgUtility.NewDefaultError("JWT generation failed", -2, 200)

		return "", oGenerateError
	}

	return sAuthorization, nil
}
