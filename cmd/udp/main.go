package udp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	pkgUtility "example/pkg/utility"
)

var Command = &cobra.Command{
	Use:   "udp",
	Short: "啟動 UDP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oLogger := pkgUtility.Logger(pkgUtility.Udp)

		sAddress := ":" + bootstrap.CONFIG.SERVICES.UDP.PORT
		oAddr, err := net.ResolveUDPAddr("udp", sAddress)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("解析 UDP 位址失敗", oErrorField)
		}

		oConn, err := net.ListenUDP("udp", oAddr)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("監聽 UDP 失敗", oErrorField)
		}
		defer oConn.Close()

		// 收到中斷/終止訊號時 oCtx 會被取消，主動關掉 conn 讓 ReadFromUDP 中斷返回，
		// 不是靠 process 被系統強制殺掉才釋放 port。
		go func() {
			<-oCtx.Done()
			oConn.Close()
		}()

		oLogger.Info("啟動 UDP 服務。 port: " + bootstrap.CONFIG.SERVICES.UDP.PORT)

		aBuf := make([]byte, 1024)

		for {
			iCount, oRemoteAddr, err := oConn.ReadFromUDP(aBuf)
			if err != nil {
				// oCtx 取消（優雅關機）就正常退出，否則是真的讀取錯誤，直接結束。
				select {
				case <-oCtx.Done():
					return
				default:
					oErrorField := zap.Error(err)
					oLogger.Fatal("讀取 UDP 失敗", oErrorField)
				}
			}

			sRemoteAddress := oRemoteAddr.String()
			sData := string(aBuf[:iCount])
			oRemoteField := zap.String("remote", sRemoteAddress)
			oDataField := zap.String("data", sData)
			oLogger.Info("收到 UDP 封包", oRemoteField, oDataField)

			aOkResponse := []byte("OK")
			oConn.WriteToUDP(aOkResponse, oRemoteAddr)
		}
	},
}
