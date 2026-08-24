package register

import (
	container "example/container"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebsocketInit 只回傳 method -> handler 對照表，不碰 upgrade/ping/dispatch 這些
// websocket 協定細節，那些是 cmd/websocket.go 用 pkg.WebsocketRouter 組裝的事。
func WebsocketInit(oContainer *container.WebsocketContainer) http.HandlerFunc {

	return func(oWriter http.ResponseWriter, oRequest *http.Request) {
		var oUpgrader = websocket.Upgrader{
			CheckOrigin: func(oRequest *http.Request) bool {
				return true
			},
		}

		oConn, oErr := oUpgrader.Upgrade(oWriter, oRequest, nil)

		if oErr != nil {
			log.Println(oErr)
			return
		}

		defer oConn.Close()

		// 讀取 client 訊息
		go func() {

			for {

				iMessageType, aMsg, oErr := oConn.ReadMessage()

				fmt.Println(iMessageType)
				fmt.Println(aMsg)

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

	}
}
