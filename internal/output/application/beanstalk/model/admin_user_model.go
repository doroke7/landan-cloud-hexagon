package outputApplicationBeanstalkModel

import (
	"encoding/json"
	"errors"
	"time"

	domain "example/internal/domain"
	outputApplicationBeanstalk "example/internal/output/application/beanstalk"
)

type AdminUserModel struct {
	*outputApplicationBeanstalk.AbstractBeanstalk
}

func NewAdminUserModel(oAbstractBeanstalk *outputApplicationBeanstalk.AbstractBeanstalk) (*AdminUserModel, error) {
	return &AdminUserModel{
		AbstractBeanstalk: oAbstractBeanstalk,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	_, err = oSelf.AbstractBeanstalk.Tube("/Queue/AdminUser.AddOne").Put(aByteBody, 1, 0, 10*time.Second)

	return err
}

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by beanstalk")
}

// Close 是空實作：Conn 現在是從 AbstractBeanstalk 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserModel) Close() error {
	return nil
}
