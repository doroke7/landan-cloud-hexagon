package cmd

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	container "example/container"
	registerFacade "example/internal/register/facade"
	pkgUtility "example/pkg/utility"
)

var oFacadeCommand = &cobra.Command{
	Use:   "facade",
	Short: "啟動 Facade 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oLogger := pkgUtility.Logger(pkgUtility.Facade)

		oContainer, err := container.InitFacadeContainer(oCtx)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("初始化 facade container 失敗", oErrorField)
		}
		oFacadeServer := registerFacade.Init(oContainer)

		sAddress := ":" + bootstrap.CONFIG.SERVICES.FACADE.PORT
		oListener, err := net.Listen("tcp", sAddress)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("監聽 FACADE port 失敗", oErrorField)
		}
		oLogger.Info("啟動 FACADE 服務。 port: " + bootstrap.CONFIG.SERVICES.FACADE.PORT)

		go func() {
			<-oCtx.Done()
			oFacadeServer.GracefulStop()
		}()

		if err := oFacadeServer.Serve(oListener); err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("FACADE server 異常結束", oErrorField)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oFacadeCommand)
}
