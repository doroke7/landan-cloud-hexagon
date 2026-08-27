package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	container "example/container"
	registerWebsocket "example/internal/register/websocket"
	pkgUtility "example/pkg/utility"
)

var oWebsocketCommand = &cobra.Command{
	Use:   "websocket",
	Short: "啟動 Websocket 服務",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitWebsocketContainer(ctx)
		if err != nil {
			log.Fatal(err)
		}
		// Websocket 才是主要關心的 服務， 所以應該 從 register 取出 websocket 套件
		oMux := registerWebsocket.Init(oContainer)

		oWebsocketServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT,
			Handler: oMux,
		}
		pkgUtility.Logger(pkgUtility.Default).Info("啟動 WEBSOCKET 服務。 port: " + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT)

		oWebsocketServer.ListenAndServe()

	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oWebsocketCommand)
}
