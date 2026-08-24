package register

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

// SocketioInit 組裝 event -> handler 對照表，回傳掛好 route 的 *http.ServeMux，
// serve/close 這些連線生命週期管理的事交給 cmd/socketio.go 做，跟 WebsocketInit 是同一套慣例。
func SocketioInit() (*socketio.Server, *http.ServeMux) {

	oServer := socketio.NewServer(nil)

	oServer.OnConnect("/", func(oConn socketio.Conn) error {

		oConn.SetContext("")
		log.Println("connected:", oConn.ID())

		return nil

	})

	oServer.OnEvent("/", "message", func(oConn socketio.Conn, sMsg string) string {

		// echo 回去
		return sMsg

	})

	oServer.OnEvent("/", "broadcast", func(oConn socketio.Conn, sMsg string) {

		oServer.BroadcastToNamespace("/", "broadcast", sMsg)

	})

	oServer.OnError("/", func(oConn socketio.Conn, oErr error) {
		log.Println("error:", oErr)
	})

	oServer.OnDisconnect("/", func(oConn socketio.Conn, sReason string) {
		log.Println("disconnected:", oConn.ID(), sReason)
	})

	oMux := http.NewServeMux()
	oMux.Handle("/socket.io/", oServer)

	return oServer, oMux

}
