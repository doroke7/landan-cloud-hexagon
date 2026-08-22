package bootstrap

import (
	"github.com/IBM/sarama"
)

func NewKafka() (sarama.Client, error) {
	oConfig := sarama.NewConfig()

	if sVersion := CONFIG.KAFKA.VERSION; sVersion != "" {
		oVersion, err := sarama.ParseKafkaVersion(sVersion)
		if err != nil {
			return nil, err
		}
		oConfig.Version = oVersion
	}

	if CONFIG.KAFKA.USER != "" {
		oConfig.Net.SASL.Enable = true
		oConfig.Net.SASL.User = CONFIG.KAFKA.USER
		oConfig.Net.SASL.Password = CONFIG.KAFKA.PASS
	}

	return sarama.NewClient(CONFIG.KAFKA.BROKERS, oConfig)
}
