package outputApplicationEmqxEvent

import (
	"encoding/json"

	domain "example/internal/domain"
	outputApplicationEmqx "example/internal/output/application/emqx"
)

type AdminUserEvent struct {
	*outputApplicationEmqx.AbstractEmqx
}

func NewAdminUserEvent(oAbstractEmqx *outputApplicationEmqx.AbstractEmqx) (*AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractEmqx: oAbstractEmqx,
	}, nil
}

func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	oToken := oSelf.Client.Publish("/Queue/AdminUser.AddOne", 0, false, aByteBody)
	oToken.Wait()

	oErr := oToken.Error()

	return oErr
}

func (oSelf *AdminUserEvent) Close() error {
	return nil
}
