package bootstrap

import (
	"strings"

	"github.com/IBM/sarama"
	"github.com/charmbracelet/log"
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

	// sarama.NewSyncProducerFromClient 要求 Producer.Return.Successes 是 true，
	// 這裡先開起來，讓 output adapter 可以直接拿這個 client 建 SyncProducer。
	oConfig.Producer.Return.Successes = true

	oClient, oErr := sarama.NewClient(CONFIG.REDPANDA.BROKERS, oConfig)

	log.Info("[INFO] REDPANDA 連線完成. ", strings.Join(CONFIG.REDPANDA.BROKERS, ","))

	return oClient, oErr
}
