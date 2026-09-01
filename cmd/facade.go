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
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitFacadeContainer(ctx)
		if err != nil {
			pkgUtility.Logger(pkgUtility.Default).Fatal("初始化 facade container 失敗", zap.Error(err))
		}
		oFacadeServer := registerFacade.Init(oContainer)

		oListener, err := net.Listen("tcp", ":"+bootstrap.CONFIG.SERVICES.FACADE.PORT)
		if err != nil {
			pkgUtility.Logger(pkgUtility.Default).Fatal("監聽 FACADE port 失敗", zap.Error(err))
		}
		pkgUtility.Logger(pkgUtility.Default).Info("啟動 FACADE 服務。 port: " + bootstrap.CONFIG.SERVICES.FACADE.PORT)

		go func() {
			<-ctx.Done()
			oFacadeServer.GracefulStop()
		}()

		if err := oFacadeServer.Serve(oListener); err != nil {
			pkgUtility.Logger(pkgUtility.Default).Fatal("FACADE server 異常結束", zap.Error(err))
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oFacadeCommand)
}
