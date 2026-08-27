package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	registerSocketio "example/internal/register/socketio"
	pkgUtility "example/pkg/utility"
)

/*
	event: connect,                                                reply ✅
	       disconnect,                                              reply ❌
		   message,                                                 reply ✅ (echo)
		   broadcast,                                                reply ✅ + broadcast ✅

    server -> client:
	       server ping (每 5 秒推播一次)                              reply ✅
*/

var oSocketioCommand = &cobra.Command{
	Use:   "socketio",
	Short: "啟動 socketio 服務",
	Run: func(cmd *cobra.Command, args []string) {

		oServer, oMux := registerSocketio.Init()

		defer oServer.Close()

		pkgUtility.Logger(pkgUtility.Default).Info("啟動 SOCKETIO 服務。 port: " + bootstrap.CONFIG.SERVICES.SOCKETIO.PORT)

		oSocketioServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.SOCKETIO.PORT,
			Handler: oMux,
		}

		log.Fatal(oSocketioServer.ListenAndServe())

	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oSocketioCommand)
}
