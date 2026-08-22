package activemq

import (
	"encoding/json"
	"errors"

	domain "example/internal/domain"
	outputApplicationActivemq "example/internal/output/application/activemq"
)

type AdminUserModel struct {
	*outputApplicationActivemq.AbstractActivemq
}

func NewAdminUserModel(oAbstractActivemq *outputApplicationActivemq.AbstractActivemq) (*AdminUserModel, error) {
	return &AdminUserModel{
		AbstractActivemq: oAbstractActivemq,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	return oSelf.Conn.Send("/queue/AdminUser.AddOne", "application/json", aByteBody)
}

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by activemq")
}

// Close 是空實作：Conn 現在是從 AbstractActivemq 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserModel) Close() error {
	return nil
}
