package bootstrap

import (
	"strings"

	"github.com/IBM/sarama"
	"github.com/charmbracelet/log"
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

	// sarama.NewSyncProducerFromClient 要求 Producer.Return.Successes 是 true，
	// 這裡先開起來，讓 output adapter 可以直接拿這個 client 建 SyncProducer。
	oConfig.Producer.Return.Successes = true

	oClient, oErr := sarama.NewClient(CONFIG.KAFKA.BROKERS, oConfig)

	log.Info("[INFO] KAFKA 連線完成.", "addr", strings.Join(CONFIG.KAFKA.BROKERS, ","))

	return oClient, oErr
}
