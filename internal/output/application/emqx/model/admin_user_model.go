package emqx

import (
	"encoding/json"
	"errors"

	domain "example/internal/domain"
	outputApplicationEmqx "example/internal/output/application/emqx"
)

type AdminUserModel struct {
	*outputApplicationEmqx.AbstractEmqx
}

func NewAdminUserModel(oAbstractEmqx *outputApplicationEmqx.AbstractEmqx) (*AdminUserModel, error) {
	return &AdminUserModel{
		AbstractEmqx: oAbstractEmqx,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	oToken := oSelf.Client.Publish("/Queue/AdminUser.AddOne", 0, false, aByteBody)
	oToken.Wait()

	oErr := oToken.Error()

	return oErr
}

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by emqx")
}

// Close 是空實作：Client 現在是從 AbstractEmqx 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserModel) Close() error {
	return nil
}
