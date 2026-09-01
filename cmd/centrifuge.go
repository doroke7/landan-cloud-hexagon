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
			pkgUtility.Logger(pkgUtility.Centrifuge).Fatal("初始化 centrifuge container 失敗", zap.Error(oErr))
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

		pkgUtility.Logger(pkgUtility.Centrifuge).Info("啟動 CENTRIFUGE 服務。 port: " + bootstrap.CONFIG.SERVICES.CENTRIFUGE.PORT)

		if oErr := oCentrifugeServer.ListenAndServe(); oErr != nil && oErr != http.ErrServerClosed {
			pkgUtility.Logger(pkgUtility.Centrifuge).Fatal("CENTRIFUGE server 異常結束", zap.Error(oErr))
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oCentrifugeCommand)
}
