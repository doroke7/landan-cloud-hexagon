package bootstrap

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Mosquitto 是 MQTT broker，這裡用 paho.mqtt.golang 這個通用 MQTT client 連線。
func NewMosquitto() (mqtt.Client, error) {
	sBroker := fmt.Sprintf("tcp://%s:%s", CONFIG.MOSQUITTO.HOST, CONFIG.MOSQUITTO.PORT)

	oOptions := mqtt.NewClientOptions().
		AddBroker(sBroker).
		SetConnectTimeout(time.Duration(CONFIG.MOSQUITTO.TIMEOUT) * time.Millisecond)

	if CONFIG.MOSQUITTO.USER != "" {
		oOptions.SetUsername(CONFIG.MOSQUITTO.USER)
		oOptions.SetPassword(CONFIG.MOSQUITTO.PASS)
	}

	oClient := mqtt.NewClient(oOptions)

	if oToken := oClient.Connect(); oToken.WaitTimeout(time.Duration(CONFIG.MOSQUITTO.TIMEOUT)*time.Millisecond) && oToken.Error() != nil {
		return nil, oToken.Error()
	}

	return oClient, nil
}
