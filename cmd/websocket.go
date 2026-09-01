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
	container "example/container"
	registerWebsocket "example/internal/register/websocket"
	pkgUtility "example/pkg/utility"
)

var oWebsocketCommand = &cobra.Command{
	Use:   "websocket",
	Short: "啟動 Websocket 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitWebsocketContainer(oCtx)
		if err != nil {
			pkgUtility.Logger(pkgUtility.Websocket).Fatal("初始化 websocket container 失敗", zap.Error(err))
		}
		// Websocket 才是主要關心的 服務， 所以應該 從 register 取出 websocket 套件
		oMux := registerWebsocket.Init(oContainer)

		oWebsocketServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT,
			Handler: oMux,
		}
		pkgUtility.Logger(pkgUtility.Websocket).Info("啟動 WEBSOCKET 服務。 port: " + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT)

		go func() {
			<-oCtx.Done()
			oWebsocketServer.Shutdown(context.Background())
		}()

		if err := oWebsocketServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkgUtility.Logger(pkgUtility.Websocket).Fatal("WEBSOCKET server 異常結束", zap.Error(err))
		}

	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oWebsocketCommand)
}
