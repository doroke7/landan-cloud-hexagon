package bootstrap

import (
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/charmbracelet/log"
)

func NewPulsar() (pulsar.Client, error) {
	sURL := fmt.Sprintf("pulsar://%s:%s", CONFIG.PULSAR.HOST, CONFIG.PULSAR.PORT)

	oClient, oErr := pulsar.NewClient(pulsar.ClientOptions{
		URL:               sURL,
		ConnectionTimeout: time.Duration(CONFIG.PULSAR.TIMEOUT) * time.Millisecond,
	})

	log.Info("[INFO] PULSAR 連線完成.", "addr", sURL)

	return oClient, oErr
}
