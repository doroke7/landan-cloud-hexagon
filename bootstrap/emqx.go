package bootstrap

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// EMQX 也是 MQTT broker，跟 mosquitto.go 一樣用 paho.mqtt.golang 連線，
// 只是連到不同的 broker 端點與帳密設定。
func NewEmqx() (mqtt.Client, error) {
	sBroker := fmt.Sprintf("tcp://%s:%s", CONFIG.EMQX.HOST, CONFIG.EMQX.PORT)

	oOptions := mqtt.NewClientOptions().
		AddBroker(sBroker).
		SetConnectTimeout(time.Duration(CONFIG.EMQX.TIMEOUT) * time.Millisecond)

	if CONFIG.EMQX.USERNAME != "" {
		oOptions.SetUsername(CONFIG.EMQX.USERNAME)
		oOptions.SetPassword(CONFIG.EMQX.PASSWORD)
	}

	oClient := mqtt.NewClient(oOptions)

	if oToken := oClient.Connect(); oToken.WaitTimeout(time.Duration(CONFIG.EMQX.TIMEOUT)*time.Millisecond) && oToken.Error() != nil {
		return nil, oToken.Error()
	}

	log.Info("[INFO] EMQX 連線完成.", "addr", sBroker)

	return oClient, nil
}
