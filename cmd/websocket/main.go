package websocket

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	container "example/container"
	registerWebsocket "example/internal/register/websocket"
	pkgUtility "example/pkg/utility"
)

var Command = &cobra.Command{
	Use:   "websocket",
	Short: "啟動 Websocket 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oLogger := pkgUtility.Logger(pkgUtility.Websocket)

		oContainer, err := container.InitWebsocketContainer(oCtx)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("初始化 websocket container 失敗", oErrorField)
		}
		// Websocket 才是主要關心的 服務， 所以應該 從 register 取出 websocket 套件
		oMux := registerWebsocket.Init(oContainer)

		oWebsocketServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT,
			Handler: oMux,
		}
		oLogger.Info("啟動 WEBSOCKET 服務。 port: " + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT)

		go func() {
			<-oCtx.Done()
			oBackgroundContext := context.Background()
			oWebsocketServer.Shutdown(oBackgroundContext)
		}()

		if err := oWebsocketServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			oErrorField := zap.Error(err)
			oLogger.Fatal("WEBSOCKET server 異常結束", oErrorField)
		}

	},
}
