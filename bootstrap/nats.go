package bootstrap

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func NewNats() (*nats.Conn, error) {
	sURL := fmt.Sprintf("nats://%s:%s", CONFIG.NATS.HOST, CONFIG.NATS.PORT)

	aOptions := []nats.Option{
		nats.Timeout(time.Duration(CONFIG.NATS.TIMEOUT) * time.Millisecond),
	}

	if CONFIG.NATS.USER != "" {
		aOptions = append(aOptions, nats.UserInfo(CONFIG.NATS.USER, CONFIG.NATS.PASS))
	}

	return nats.Connect(sURL, aOptions...)
}
