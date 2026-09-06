package outputApplicationEmqxEvent

import (
	"encoding/json"

	domain "example/internal/domain"
	outputPortAnyEvent "example/internal/output/port/any/event"
)

type AdminUserEvent struct {
	*AbstractEvent
}

func NewAdminUserEvent(oAbstractEvent *AbstractEvent) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractEvent: oAbstractEvent,
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
