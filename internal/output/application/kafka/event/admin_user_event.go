package outputApplicationKafkaEvent

import (
	"encoding/json"

	"github.com/IBM/sarama"

	domain "example/internal/domain"
	outputApplicationKafka "example/internal/output/application/kafka"
	outputPortAnyEvent "example/internal/output/port/any/event"
)

type AdminUserEvent struct {
	*outputApplicationKafka.AbstractKafka
}

func NewAdminUserEvent(oAbstractKafka *outputApplicationKafka.AbstractKafka) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractKafka: oAbstractKafka,
	}, nil
}

func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
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

// Close 是空實作：Client/Producer 現在都是從 AbstractKafka 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserEvent) Close() error {
	return nil
}
