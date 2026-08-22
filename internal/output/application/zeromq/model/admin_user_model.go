package zeromq

import (
	"encoding/json"
	"errors"

	"github.com/go-zeromq/zmq4"

	domain "example/internal/domain"
	outputApplicationZeromq "example/internal/output/application/zeromq"
)

type AdminUserModel struct {
	*outputApplicationZeromq.AbstractZeromq
}

func NewAdminUserModel(oAbstractZeromq *outputApplicationZeromq.AbstractZeromq) (*AdminUserModel, error) {
	return &AdminUserModel{
		AbstractZeromq: oAbstractZeromq,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	return oSelf.Socket.SendMulti(zmq4.NewMsgFrom(
		[]byte("/Queue/AdminUser.AddOne"),
		aByteBody,
	))
}

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by zeromq")
}

// Close 是空實作：Socket 現在是從 AbstractZeromq 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserModel) Close() error {
	return nil
}
