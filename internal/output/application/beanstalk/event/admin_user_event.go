package outputApplicationBeanstalkEvent

import (
	"encoding/json"
	"time"

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

	_, err = oSelf.AbstractBeanstalk.Tube("/Queue/AdminUser.AddOne").Put(aByteBody, 1, 0, 10*time.Second)

	return err
}

func (oSelf *AdminUserEvent) Close() error {
	return nil
}
