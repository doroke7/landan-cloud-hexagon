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
	registerCentrifuge "example/internal/register/centrifuge"
	pkgUtility "example/pkg/utility"
)

var oCentrifugeCommand = &cobra.Command{
	Use:   "centrifuge",
	Short: "啟動 centrifuge 服務",
	Run: func(cmd *cobra.Command, args []string) {

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, oErr := container.InitCentrifugeContainer(ctx)
		if oErr != nil {
			log.Fatal(oErr)
		}
		defer oContainer.Nats.Close()

		oNode, oHandler := registerCentrifuge.Init(oContainer)

		oCentrifugeServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.CENTRIFUGE.PORT,
			Handler: oHandler,
		}

		go func() {
			<-ctx.Done()
			oNode.Shutdown(context.Background())
			oCentrifugeServer.Shutdown(context.Background())
		}()

		pkgUtility.Logger(pkgUtility.Default).Info("啟動 CENTRIFUGE 服務。 port: " + bootstrap.CONFIG.SERVICES.CENTRIFUGE.PORT)

		if oErr := oCentrifugeServer.ListenAndServe(); oErr != nil && oErr != http.ErrServerClosed {
			log.Fatal(oErr)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oCentrifugeCommand)
}
