package nats

import (
	"encoding/json"
	"errors"

	domain "example/internal/domain"
	outputApplicationNats "example/internal/output/application/nats"
)

type AdminUserModel struct {
	*outputApplicationNats.AbstractNats
}

func NewAdminUserModel(oAbstractNats *outputApplicationNats.AbstractNats) (*AdminUserModel, error) {
	return &AdminUserModel{
		AbstractNats: oAbstractNats,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	return oSelf.Conn.Publish("AdminUser.AddOne", aByteBody)
}

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by nats")
}

// Close 是空實作：Conn 現在是從 AbstractNats 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserModel) Close() error {
	return nil
}
