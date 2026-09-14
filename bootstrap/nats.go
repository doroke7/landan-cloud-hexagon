package bootstrap

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
)

func NewNats() (*nats.Conn, error) {
	sURL := fmt.Sprintf("nats://%s:%s", CONFIG.NATS.HOST, CONFIG.NATS.PORT)

	aOptions := []nats.Option{
		nats.Timeout(time.Duration(CONFIG.NATS.TIMEOUT) * time.Millisecond),
	}

	if CONFIG.NATS.USERNAME != "" {
		oUserInfoOption := nats.UserInfo(CONFIG.NATS.USERNAME, CONFIG.NATS.PASSWORD)
		aOptions = append(aOptions, oUserInfoOption)
	}

	oConn, oErr := nats.Connect(sURL, aOptions...)

	log.Info("[INFO] NATS 連線完成.", "addr", CONFIG.NATS.HOST+":"+CONFIG.NATS.PORT)

	return oConn, oErr
}
