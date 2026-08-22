package redpanda

import (
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"

	domain "example/internal/domain"
	outputApplicationRedpanda "example/internal/output/application/redpanda"
)

type AdminUserModel struct {
	*outputApplicationRedpanda.AbstractRedpanda
}

func NewAdminUserModel(oAbstractRedpanda *outputApplicationRedpanda.AbstractRedpanda) (*AdminUserModel, error) {
	return &AdminUserModel{
		AbstractRedpanda: oAbstractRedpanda,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	_, _, err = oSelf.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "/Queue/AdminUser.AddOne",
		Value: sarama.ByteEncoder(aByteBody),
	})

	return err
}

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by redpanda")
}

// Close 是空實作：Client/Producer 現在都是從 AbstractRedpanda 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserModel) Close() error {
	return nil
}
