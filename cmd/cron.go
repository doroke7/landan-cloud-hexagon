package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	container "example/container"
	registerCron "example/internal/register/cron"
	pkgUtility "example/pkg/utility"
)

var oCronCommand = &cobra.Command{
	Use:   "cron",
	Short: "啟動排程服務",
	Run: func(cmd *cobra.Command, args []string) {
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitCronContainer(oCtx)
		if err != nil {
			pkgUtility.Logger(pkgUtility.Cron).Fatal("初始化 cron container 失敗", zap.Error(err))
		}

		oCron := registerCron.Init(oContainer)

		pkgUtility.Logger(pkgUtility.Cron).Info("啟動 CRON 服務。")

		oCron.Start()

		<-oCtx.Done()
		oCron.Stop()
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oCronCommand)
}
