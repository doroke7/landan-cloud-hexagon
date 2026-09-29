package cron

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

var Command = &cobra.Command{
	Use:   "cron",
	Short: "啟動排程服務",
	Run: func(cmd *cobra.Command, args []string) {
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oLogger := pkgUtility.Logger(pkgUtility.Cron)

		oContainer, err := container.InitCronContainer(oCtx)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("初始化 cron container 失敗", oErrorField)
		}

		oCron := registerCron.Init(oContainer)

		oLogger.Info("啟動 CRON 服務。")

		oCron.Start()

		<-oCtx.Done()
		oCron.Stop()
	},
}
