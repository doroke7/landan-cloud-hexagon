package centrifuge

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

var Command = &cobra.Command{
	Use:   "centrifuge",
	Short: "啟動 centrifuge 服務",
	Run: func(cmd *cobra.Command, args []string) {

		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oLogger := pkgUtility.Logger(pkgUtility.Centrifuge)

		oContainer, oErr := container.InitCentrifugeContainer(oCtx)
		if oErr != nil {
			oErrorField := zap.Error(oErr)
			oLogger.Fatal("初始化 centrifuge container 失敗", oErrorField)
		}
		defer oContainer.Nats.Close()

		oNode, oHandler := registerCentrifuge.Init(oContainer)

		oCentrifugeServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.CENTRIFUGE.PORT,
			Handler: oHandler,
		}

		go func() {
			<-oCtx.Done()
			oBackgroundContext := context.Background()
			oNode.Shutdown(oBackgroundContext)
			oCentrifugeServer.Shutdown(oBackgroundContext)
		}()

		oLogger.Info("啟動 CENTRIFUGE 服務。 port: " + bootstrap.CONFIG.SERVICES.CENTRIFUGE.PORT)

		if oErr := oCentrifugeServer.ListenAndServe(); oErr != nil && oErr != http.ErrServerClosed {
			oErrorField := zap.Error(oErr)
			oLogger.Fatal("CENTRIFUGE server 異常結束", oErrorField)
		}
	},
}
