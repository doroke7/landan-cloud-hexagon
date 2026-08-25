package register

import (
	"log"
	"net/http"

	socketio "github.com/doquangtan/socketio/v4"
	"github.com/rs/cors"
)

func SocketioInit() (*socketio.Io, *http.ServeMux) {

	oCors := cors.New(cors.Options{
		AllowOriginFunc:  func(sOrigin string) bool { return true },
		AllowCredentials: true,
	})

	oServer := socketio.New()

	oServer.OnConnection(func(oSocket *socketio.Socket) {

		log.Println("connected:", oSocket.Id)

		oSocket.On("message", func(oEvent *socketio.EventPayload) {

			// echo 回去
			oSocket.Emit("message", oEvent.Data...)

		})

		oSocket.On("broadcast", func(oEvent *socketio.EventPayload) {
			log.Println("broadcast:", oEvent.Data)

			// oServer.Emit("broadcast", oEvent.Data...)

		})

		oSocket.On("disconnect", func(oEvent *socketio.EventPayload) {
			log.Println("disconnected:", oSocket.Id)
		})

	})

	oMux := http.NewServeMux()
	oMux.Handle("/socket.io/", oCors.Handler(oServer.HttpHandler()))

	return oServer, oMux

}
