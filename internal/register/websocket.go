package register

import (
	container "example/container"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebsocketInit 組裝 upgrade/echo 這些 websocket 協定細節，回傳掛好 route 的 *http.ServeMux，
// serve 的事交給 cmd/websocket.go 做，跟 SocketioInit 是同一套慣例。
func WebsocketInit(oContainer *container.WebsocketContainer) *http.ServeMux {

	oMux := http.NewServeMux()

	oMux.HandleFunc("/ws", func(oWriter http.ResponseWriter, oRequest *http.Request) {
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

				if oErr != nil {
					fmt.Println("iMessageType", iMessageType)
					fmt.Println("aMsg=", aMsg)
					fmt.Println("oErr=", oErr)

					return
				}

				// echo 回去
				oErr = oConn.WriteMessage(iMessageType, aMsg)

				if oErr != nil {
					return
				}

			}

		}()

	})

	return oMux

}
