package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	container "example/container"
	register "example/internal/register"
	pkg "example/pkg"
)

var oHttpCommand = &cobra.Command{
	Use:   "http",
	Short: "啟動 Gin HTTP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitHttpContainer(ctx)
		if err != nil {
			log.Fatal(err)
		}
		oGin := gin.Default()

		oEngine := register.HttpInit(oGin, oContainer)

		// gin.Engine.Run() 內部自己建立 http.Server、拿不到參考做 Shutdown，
		// 改成自己組 http.Server，收到中斷/終止訊號時主動 Shutdown，
		// 讓 ListenAndServe() 正常返回並釋放 port，不是靠 process 被強制殺掉才釋放。
		oHttpServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.HTTP.PORT,
			Handler: oEngine,
		}

		go func() {
			<-ctx.Done()
			oHttpServer.Shutdown(context.Background())
		}()

		pkg.Logger(pkg.Default).Info("啟動 RESOURCE 服務。 port: " + bootstrap.CONFIG.SERVICES.HTTP.PORT)

		if err := oHttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oHttpCommand)
}
