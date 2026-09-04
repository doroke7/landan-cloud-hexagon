package helper

/*
    IMPORTANT: 統一改寫為 helper 再注入框架
	優點：寫法一致結構性強
*/

import (
	"github.com/go-playground/validator/v10"
)

type ValidatorHelper struct {
	*AbstractHelper
	validator *validator.Validate
}

func NewValiatorHelper(oAbstractHelper *AbstractHelper) *ValidatorHelper {
	return &ValidatorHelper{
		AbstractHelper: oAbstractHelper,
		validator:      validator.New(),
	}
}

func (oSelf *ValidatorHelper) Struct(oStruct any) error {
	oError := oSelf.validator.Struct(oStruct)
	return oError
}
