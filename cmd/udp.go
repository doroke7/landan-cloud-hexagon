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
	pkgUtility "example/pkg/utility"
)

var oUdpCommand = &cobra.Command{
	Use:   "udp",
	Short: "啟動 UDP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oAddr, err := net.ResolveUDPAddr("udp", ":"+bootstrap.CONFIG.SERVICES.UDP.PORT)
		if err != nil {
			pkgUtility.Logger(pkgUtility.Udp).Fatal("解析 UDP 位址失敗", zap.Error(err))
		}

		oConn, err := net.ListenUDP("udp", oAddr)
		if err != nil {
			pkgUtility.Logger(pkgUtility.Udp).Fatal("監聽 UDP 失敗", zap.Error(err))
		}
		defer oConn.Close()

		// 收到中斷/終止訊號時 oCtx 會被取消，主動關掉 conn 讓 ReadFromUDP 中斷返回，
		// 不是靠 process 被系統強制殺掉才釋放 port。
		go func() {
			<-oCtx.Done()
			oConn.Close()
		}()

		pkgUtility.Logger(pkgUtility.Udp).Info("啟動 UDP 服務。 port: " + bootstrap.CONFIG.SERVICES.UDP.PORT)

		aBuf := make([]byte, 1024)

		for {
			iCount, oRemoteAddr, err := oConn.ReadFromUDP(aBuf)
			if err != nil {
				// oCtx 取消（優雅關機）就正常退出，否則是真的讀取錯誤，直接結束。
				select {
				case <-oCtx.Done():
					return
				default:
					pkgUtility.Logger(pkgUtility.Udp).Fatal("讀取 UDP 失敗", zap.Error(err))
				}
			}

			pkgUtility.Logger(pkgUtility.Udp).Info("收到 UDP 封包",
				zap.String("remote", oRemoteAddr.String()),
				zap.String("data", string(aBuf[:iCount])),
			)

			oConn.WriteToUDP([]byte("OK"), oRemoteAddr)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oUdpCommand)
}
