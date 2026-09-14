package helper

/*
    IMPORTANT: 統一改寫為 helper 再注入框架
	優點：寫法一致結構性強
*/

import (
	pkgUtility "example/pkg/utility"

	"github.com/go-playground/validator/v10"
)

type ValidatorHelper struct {
	*AbstractHelper
	validator *validator.Validate
}

func NewValiatorHelper(oAbstractHelper *AbstractHelper) *ValidatorHelper {
	oValidatorHelper := &ValidatorHelper{
		AbstractHelper: oAbstractHelper,
		validator:      validator.New(),
	}

	return oValidatorHelper
}

func (oSelf *ValidatorHelper) Struct(oStruct any) error {
	oError := oSelf.validator.Struct(oStruct)
	return oError
}

func (oSelf *ValidatorHelper) Valiate(oStruct any) error {
	oError := oSelf.validator.Struct(oStruct)

	if oError != nil {
		oErrors := oError.(validator.ValidationErrors)
		oError := oErrors[0]

		sField := oError.Field()
		sTag := oError.Tag()
		sParam := oError.Param()
		sMessage := sField + " requires " + sTag + " " + sParam
		oDefaultError := pkgUtility.NewDefaultError(sMessage, -1, 200)

		return oDefaultError

	}

	return oError
}
