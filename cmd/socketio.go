package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

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
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oServer, oMux := registerSocketio.Init()

		defer oServer.Close()

		oSocketioServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.SOCKETIO.PORT,
			Handler: oMux,
		}

		go func() {
			<-oCtx.Done()
			oSocketioServer.Shutdown(context.Background())
		}()

		pkgUtility.Logger(pkgUtility.Socketio).Info("啟動 SOCKETIO 服務。 port: " + bootstrap.CONFIG.SERVICES.SOCKETIO.PORT)

		if err := oSocketioServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkgUtility.Logger(pkgUtility.Socketio).Fatal("SOCKETIO server 異常結束", zap.Error(err))
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oSocketioCommand)
}
