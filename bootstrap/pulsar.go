package bootstrap

import (
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func NewPulsar() (pulsar.Client, error) {
	sURL := fmt.Sprintf("pulsar://%s:%s", CONFIG.PULSAR.HOST, CONFIG.PULSAR.PORT)

	return pulsar.NewClient(pulsar.ClientOptions{
		URL:               sURL,
		ConnectionTimeout: time.Duration(CONFIG.PULSAR.TIMEOUT) * time.Millisecond,
	})
}
