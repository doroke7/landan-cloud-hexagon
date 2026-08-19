package bootstrap

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	utility "example/internal/utility"
)

func NewResource(oContext context.Context) *grpc.ClientConn {

	aHosts := CONFIG.CLIENTS.RESOURCE.HOSTS
	aPorts := CONFIG.CLIENTS.RESOURCE.PORTS

	aAddrs := make([]resolver.Address, len(aHosts))
	for iIndex, sHost := range aHosts {
		sPort := aPorts[0]
		if iIndex < len(aPorts) {
			sPort = aPorts[iIndex]
		}
		aAddrs[iIndex] = resolver.Address{Addr: net.JoinHostPort(sHost, sPort)}
	}

	// manual.Resolver 直接把固定位址塞給 grpc.WithResolvers，只作用於這個 client，
	// 不用再手寫 resolver.Builder / resolver.Resolver，也不需要 resolver.Register 佔用全域 scheme。
	oResolverBuilder := manual.NewBuilderWithScheme("resource-static")
	oResolverBuilder.InitialState(resolver.State{Addresses: aAddrs})

	// gRPC 本身支持 多路復用， 不建議做連結池
	oConnection, err := grpc.NewClient(oResolverBuilder.Scheme()+":///resource", grpc.WithResolvers(oResolverBuilder), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		sUser := CONFIG.CLIENTS.RESOURCE.USER
		sPassword := CONFIG.CLIENTS.RESOURCE.PASSWORD

		sAuthorization := "Basic " + utility.Base64Encode(
			sUser+":"+sPassword,
		)
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", sAuthorization))
		return invoker(ctx, method, req, reply, cc, opts...)
	}),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1.0 * time.Second, // 第一次斷線後，等 1.0 秒再嘗試重連
				Multiplier: 1.6,               // 每次重連失敗，等待時間乘以 1.6 (1s -> 1.6s -> 2.56s)
				Jitter:     0.2,               // 加上 20% 的隨機抖動誤差，把大量 Client 的重連時間錯開
				MaxDelay:   10 * time.Second,  // 不管失敗幾次，最長只等 30 秒，避免時間被無限拉長
			},
			MinConnectTimeout: 3 * time.Second, // 每次嘗試建立 TCP 握手時，最少給底層 3 秒的超時時間
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second, // 每 10 秒偷偷送一次 PING 保活、防止被防火牆剪斷
			Timeout:             3 * time.Second,  // PING 出去後 3 秒內 Server 沒回應，直接判定斷線，立刻觸發上面的 Backoff 流程
			PermitWithoutStream: true,             // 關鍵：就算現在業務沒請求、沒有 Stream，也要送 PING
		}),
	)
	if err != nil {
		return nil
	}

	ticker := time.NewTicker(10 * time.Second)

	oConnection.Connect()

	// 在背景啟動常駐任務，檢查連線；oContext 取消（優雅關機）就跟著退出，
	// 不然這個 goroutine 會活到 process 被強制殺掉才釋放，跟 pkg.TcpRouter/
	// WebsocketRouter 那些長駐 goroutine 用 ctx 做優雅關機是同一套慣例。
	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-oContext.Done():
				return
			case <-ticker.C:
				state := oConnection.GetState()
				// oLogger.Info("⚠️ 檢查狀態")
				// 2. 核心判斷邏輯：
				// - 如果是 TransientFailure：代表剛斷線，Backoff 正在數鬧鐘，我們進來加速它，打破等待窗期。
				// - 如果是 Idle：代表因為無人點餐，Backoff 已經罷工了！我們必須進來重新點火。
				if state == connectivity.TransientFailure || state == connectivity.Idle {

					// 🔥 執行核心重連動作：打破 IDLE，逼迫底層立刻重新進行 TCP 握手
					oConnection.Connect()
				}
			}
		}

	}()

	return oConnection
}
