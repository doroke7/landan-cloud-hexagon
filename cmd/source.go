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
	registerSource "example/internal/register/source"
	pkgUtility "example/pkg/utility"
)

var oSourceCommand = &cobra.Command{
	Use:   "source",
	Short: "啟動 Rource 服務",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitSourceContainer(ctx)
		if err != nil {
			pkgUtility.Logger(pkgUtility.Source).Fatal("初始化 source container 失敗", zap.Error(err))
		}

		oResourceServer := registerSource.Init(oContainer)

		oListener, err := net.Listen("tcp", ":"+bootstrap.CONFIG.SERVICES.SOURCE.PORT)
		if err != nil {
			pkgUtility.Logger(pkgUtility.Source).Fatal("監聽 SOURCE port 失敗", zap.Error(err))
		}
		pkgUtility.Logger(pkgUtility.Source).Info("啟動 SOURCE 服務。 port: " + bootstrap.CONFIG.SERVICES.SOURCE.PORT)

		// 收到中斷/終止訊號時 ctx 會被取消，主動 GracefulStop，讓 gRPC 停止 accept 新連線、
		// 關掉 listener，Serve() 才會正常返回並釋放 port，
		// 不是靠 process 被系統強制殺掉才釋放。
		go func() {
			<-ctx.Done()
			oResourceServer.GracefulStop()
		}()

		if err := oResourceServer.Serve(oListener); err != nil {
			pkgUtility.Logger(pkgUtility.Source).Fatal("SOURCE server 異常結束", zap.Error(err))
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oSourceCommand)
}
