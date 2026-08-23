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
	register "example/internal/register"
	pkg "example/pkg"
)

var oWebsocketCommand = &cobra.Command{
	Use:   "websocket",
	Short: "啟動 websocket 服務",
	Run: func(cmd *cobra.Command, args []string) {

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitWebsocketContainer(ctx)
		if err != nil {
			log.Fatal(err)
		}

		oRouter := register.WebsocketInit(oContainer)
		oRouter.Serve(ctx)

		oWebsocketServer := &http.Server{
			Addr: ":" + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT,
		}

		go func() {
			<-ctx.Done()
			oWebsocketServer.Shutdown(context.Background())
		}()

		pkg.Logger(pkg.Default).Info("啟動 WEBSOCKET 服務。 port: " + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT)

		if err := oWebsocketServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oWebsocketCommand)
}
