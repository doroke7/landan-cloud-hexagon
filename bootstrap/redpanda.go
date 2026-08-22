package bootstrap

import (
	"github.com/IBM/sarama"
)

// Redpanda 跟 Kafka 走同一套 wire protocol，所以直接沿用 sarama，
// 差異只在連線資訊（CONFIG.REDPANDA）跟目標叢集是 Redpanda 而不是 Kafka。
func NewRedpanda() (sarama.Client, error) {
	oConfig := sarama.NewConfig()

	if sVersion := CONFIG.REDPANDA.VERSION; sVersion != "" {
		oVersion, err := sarama.ParseKafkaVersion(sVersion)
		if err != nil {
			return nil, err
		}
		oConfig.Version = oVersion
	}

	if CONFIG.REDPANDA.USER != "" {
		oConfig.Net.SASL.Enable = true
		oConfig.Net.SASL.User = CONFIG.REDPANDA.USER
		oConfig.Net.SASL.Password = CONFIG.REDPANDA.PASS
	}

	return sarama.NewClient(CONFIG.REDPANDA.BROKERS, oConfig)
}
