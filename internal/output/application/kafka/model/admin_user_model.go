package kafka

import (
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"

	domain "example/internal/domain"
	outputApplicationKafka "example/internal/output/application/kafka"
)

type AdminUserModel struct {
	*outputApplicationKafka.AbstractKafka
}

func NewAdminUserModel(oAbstractKafka *outputApplicationKafka.AbstractKafka) (*AdminUserModel, error) {
	return &AdminUserModel{
		AbstractKafka: oAbstractKafka,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	_, _, err = oSelf.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "AdminUser.AddOne",
		Value: sarama.ByteEncoder(aByteBody),
	})

	return err
}

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by kafka")
}

// Close 是空實作：Client/Producer 現在都是從 AbstractKafka 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserModel) Close() error {
	return nil
}
