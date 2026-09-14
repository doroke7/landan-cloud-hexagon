package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	container "example/container"
	registerHttp "example/internal/register/http"
	pkgUtility "example/pkg/utility"
)

var oHttpCommand = &cobra.Command{
	Use:   "http",
	Short: "啟動 Gin HTTP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oLogger := pkgUtility.Logger(pkgUtility.Http)

		oContainer, err := container.InitHttpContainer(oCtx)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("初始化 http container 失敗", oErrorField)
		}
		oGin := gin.Default()

		oEngine := registerHttp.Init(oGin, oContainer)

		// gin.Engine.Run() 內部自己建立 http.Server、拿不到參考做 Shutdown，
		// 改成自己組 http.Server，收到中斷/終止訊號時主動 Shutdown，
		// 讓 ListenAndServe() 正常返回並釋放 port，不是靠 process 被強制殺掉才釋放。
		oHttpServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.HTTP.PORT,
			Handler: oEngine,
		}

		go func() {
			<-oCtx.Done()
			oBackgroundContext := context.Background()
			oHttpServer.Shutdown(oBackgroundContext)
		}()

		oLogger.Info("啟動 HTTP 服務。 port: " + bootstrap.CONFIG.SERVICES.HTTP.PORT)

		if err := oHttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			oErrorField := zap.Error(err)
			oLogger.Fatal("HTTP server 異常結束", oErrorField)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oHttpCommand)
}
