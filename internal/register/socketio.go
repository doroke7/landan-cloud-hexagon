package register

import (
	"net/http"

	socketio "github.com/doquangtan/socketio/v4"
	"github.com/rs/cors"
	"go.uber.org/zap"

	pkg "example/pkg"
)

func SocketioInit() (*socketio.Io, *http.ServeMux) {

	oCors := cors.New(cors.Options{
		AllowOriginFunc:  func(sOrigin string) bool { return true },
		AllowCredentials: true,
	})

	oServer := socketio.New()

	oServer.OnConnection(func(oSocket *socketio.Socket) {

		pkg.Logger(pkg.Socketio).Info("connected", zap.String("id", oSocket.Id))

		oSocket.On("message", func(oEvent *socketio.EventPayload) {

			// echo 回去
			oSocket.Emit("message", oEvent.Data...)

		})

		oSocket.On("broadcast", func(oEvent *socketio.EventPayload) {
			pkg.Logger(pkg.Socketio).Info("broadcast", zap.Any("data", oEvent.Data))

			// oServer.Emit("broadcast", oEvent.Data...)

		})

		oSocket.On("disconnect", func(oEvent *socketio.EventPayload) {
			pkg.Logger(pkg.Socketio).Info("disconnected", zap.String("id", oSocket.Id))
		})

	})

	oMux := http.NewServeMux()
	oMux.Handle("/socket.io/", oCors.Handler(oServer.HttpHandler()))

	return oServer, oMux

}
