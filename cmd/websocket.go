package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var oUpgrader = websocket.Upgrader{
	CheckOrigin: func(oRequest *http.Request) bool {
		return true
	},
}

func fnHandler(oWriter http.ResponseWriter, oRequest *http.Request) {

	oConn, oErr := oUpgrader.Upgrade(
		oWriter,
		oRequest,
		nil,
	)

	if oErr != nil {
		log.Println(oErr)
		return
	}

	defer oConn.Close()

	// 讀取 client 訊息
	go func() {

		for {

			iMessageType, aMsg, oErr := oConn.ReadMessage()

			if oErr != nil {
				log.Println("read error:", oErr)
				return
			}

			// echo 回去
			oErr = oConn.WriteMessage(iMessageType, aMsg)

			if oErr != nil {
				return
			}

		}

	}()

	// server 主動推送
	oTicker := time.NewTicker(5 * time.Second)

	defer oTicker.Stop()

	for range oTicker.C {

		oErr := oConn.WriteMessage(websocket.TextMessage, []byte("server ping"))

		if oErr != nil {
			return
		}

	}

}

func main() {

	http.HandleFunc("/ws", fnHandler)

	http.ListenAndServe(":8080", nil)

}
