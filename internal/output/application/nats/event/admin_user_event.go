package outputApplicationNatsEvent

import (
	"encoding/json"

	domain "example/internal/domain"
	outputApplicationNats "example/internal/output/application/nats"
	outputPortAnyEvent "example/internal/output/port/any/event"
)

type AdminUserEvent struct {
	*outputApplicationNats.AbstractNats
}

func NewAdminUserEvent(oAbstractNats *outputApplicationNats.AbstractNats) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractNats: oAbstractNats,
	}, nil
}

func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	return oSelf.Conn.Publish("/Queue/AdminUser.AddOne", aByteBody)
}

func (oSelf *AdminUserEvent) Close() error {
	return nil
}
