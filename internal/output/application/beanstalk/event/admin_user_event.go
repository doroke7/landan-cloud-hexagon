package outputApplicationBeanstalkEvent

import (
	"encoding/json"
	"time"

	domain "example/internal/domain"
	outputApplicationBeanstalk "example/internal/output/application/beanstalk"
	outputPortAnyEvent "example/internal/output/port/any/event"
)

type AdminUserEvent struct {
	*outputApplicationBeanstalk.AbstractBeanstalk
}

func NewAdminUserEvent(oAbstractBeanstalk *outputApplicationBeanstalk.AbstractBeanstalk) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractBeanstalk: oAbstractBeanstalk,
	}, nil
}

func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	oTube := oSelf.AbstractBeanstalk.Tube("/Queue/AdminUser.AddOne")
	_, err = oTube.Put(aByteBody, 1, 0, 10*time.Second)

	return err
}

func (oSelf *AdminUserEvent) Close() error {
	return nil
}
