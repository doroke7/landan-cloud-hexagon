package register

import (
	container "example/container"
	"log"
	"net/http"

	pkg "example/pkg"

	"github.com/gorilla/websocket"
)

// WebsocketInit 只負責注入業務邏輯（連線/斷線 log、echo），連線生命週期機制
// （upgrade、read loop、斷線偵測）交給 pkg.WebsocketRouter，跟 SocketioInit／
// CentrifugeInit 是同一套「通訊邏輯跟業務邏輯分開」的慣例。
func WebsocketInit(oContainer *container.WebsocketContainer) *http.ServeMux {

	oEventer := pkg.NewWebsocketEventer(websocket.Upgrader{
		CheckOrigin: func(oRequest *http.Request) bool {
			return true
		},
	})

	oEventer.OnConnect(func(oConn *websocket.Conn) {
		log.Println("connected:", oConn.RemoteAddr())
	})
	oEventer.OnDisconnect(func(oConn *websocket.Conn) {
		log.Println("disconnected:", oConn.RemoteAddr())
	})
	oEventer.OnMessage(func(oConn *websocket.Conn, iMessageType int, aMsg []byte) {
		// echo 回去
		oConn.WriteMessage(iMessageType, aMsg)
	})

	oMux := http.NewServeMux()
	oMux.Handle("/ws", oEventer)

	return oMux
}
