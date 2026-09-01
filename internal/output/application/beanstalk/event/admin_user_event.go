package outputApplicationBeanstalkEvent

import (
	"encoding/json"
	"time"

	domain "example/internal/domain"
	outputApplicationBeanstalk "example/internal/output/application/beanstalk"
)

type AdminUserEvent struct {
	*outputApplicationBeanstalk.AbstractBeanstalk
}

func NewAdminUserEvent(oAbstractBeanstalk *outputApplicationBeanstalk.AbstractBeanstalk) (*AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractBeanstalk: oAbstractBeanstalk,
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
