package outputApplicationMosquittoEvent

import (
	"encoding/json"

	domain "example/internal/domain"
	outputApplicationMosquitto "example/internal/output/application/mosquitto"
)

type AdminUserEvent struct {
	*outputApplicationMosquitto.AbstractMosquitto
}

func NewAdminUserEvent(oAbstractMosquitto *outputApplicationMosquitto.AbstractMosquitto) (*AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractMosquitto: oAbstractMosquitto,
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

// Close 是空實作：Client 現在是從 AbstractMosquitto 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserEvent) Close() error {
	return nil
}
