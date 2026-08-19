package grpcclient

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sort"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

type Address struct {
	Ip     string
	Port   string
	Weight uint
}

const weightedRoundRobinName = "static_weighted_round_robin"

// IMPORTANT 重試機制：最多四次
const serviceConfigJSON = `{
	"loadBalancingConfig": [{"%s": {}}],
	"methodConfig": [{
		"name": [{}],
		"retryPolicy": {
			"maxAttempts": 4,
			"initialBackoff": "0.1s",
			"maxBackoff": "1s",
			"backoffMultiplier": 2.0,
			"retryableStatusCodes": ["UNAVAILABLE"]
		}
	}]
}`

type weightAttributeKey struct{}

func init() {
	balancer.Register(base.NewBalancerBuilder(
		weightedRoundRobinName,
		&weightedPickerBuilder{},
		base.Config{HealthCheck: true},
	))
}

type weightedPickerBuilder struct{}

func (*weightedPickerBuilder) Build(oInfo base.PickerBuildInfo) balancer.Picker {
	if len(oInfo.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	aSubConns := make([]*weightedSubConn, 0, len(oInfo.ReadySCs))
	for oSc, oScInfo := range oInfo.ReadySCs {
		iWeight, _ := oScInfo.Address.BalancerAttributes.Value(weightAttributeKey{}).(uint)
		aSubConns = append(aSubConns, &weightedSubConn{sc: oSc, weight: iWeight, addr: oScInfo.Address.Addr})
	}

	// 權重高的排最前面；權重相同時用位址排序，只是為了讓「第一個」在每次 Build
	// 都是同一台，不會因為 map 疊代順序隨機而每次啟動挑到不同台。
	sort.Slice(aSubConns, func(iIndex1 int, iIndex2 int) bool {
		if aSubConns[iIndex1].weight != aSubConns[iIndex2].weight {
			return aSubConns[iIndex1].weight > aSubConns[iIndex2].weight
		}
		return aSubConns[iIndex1].addr < aSubConns[iIndex2].addr
	})

	return &weightedPicker{subConns: aSubConns}
}

type weightedSubConn struct {
	sc     balancer.SubConn
	weight uint
	addr   string
}

// weightedPicker 固定挑權重最高的第一台，權重本身不會因為被選中而遞減；
// 只有在上一次選到的那台跑出 RPC 錯誤時，才會從其他台裡隨機挑一台頂上，
// 之後同樣固定打這台，直到它也出錯才會再換下一台。
type weightedPicker struct {
	mu       sync.Mutex
	subConns []*weightedSubConn
	current  int
}

func (oSelf *weightedPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	oSelf.mu.Lock()
	oCurrent := oSelf.subConns[oSelf.current]
	oSelf.mu.Unlock()

	return balancer.PickResult{
		SubConn: oCurrent.sc,
		Done: func(oDone balancer.DoneInfo) {
			if oDone.Err != nil {
				oSelf.failover(oCurrent, oDone.Err)
			}
		},
	}, nil
}

func (oSelf *weightedPicker) failover(oFailed *weightedSubConn, oErr error) {
	oSelf.mu.Lock()
	defer oSelf.mu.Unlock()

	if len(oSelf.subConns) <= 1 {
		log.Printf("[grpc_client] conn 失敗 addr=%s err=%v，但沒有其他台可以換", oFailed.addr, oErr)
		return
	}

	// 從「除了目前這台」以外的其他台裡隨機挑一個。
	iNext := rand.Intn(len(oSelf.subConns) - 1)
	if iNext >= oSelf.current {
		iNext++
	}
	oSelf.current = iNext

	// IMPORTANT：只要觸發 換 conn 就打日誌 或發送到 運帷報警群
	log.Printf("[grpc_client] conn 失敗 addr=%s err=%v，換 conn -> addr=%s", oFailed.addr, oErr, oSelf.subConns[iNext].addr)
}

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(aAddrs []Address) (*Client, error) {

	aResolverAddrs := make([]resolver.Address, len(aAddrs))
	for iIndex, oAddr := range aAddrs {
		aResolverAddrs[iIndex] = resolver.Address{
			Addr:               net.JoinHostPort(oAddr.Ip, oAddr.Port),
			BalancerAttributes: attributes.New(weightAttributeKey{}, oAddr.Weight),
		}
	}

	oResolverBuilder := manual.NewBuilderWithScheme("grpc-client-static")
	oResolverBuilder.InitialState(resolver.State{Addresses: aResolverAddrs})

	conn, err := grpc.NewClient(
		oResolverBuilder.Scheme()+":///",
		grpc.WithResolvers(oResolverBuilder),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{
			"loadBalancingConfig": [{"%s": {}}],
			"methodConfig": [{
				"name": [{}],
				"retryPolicy": {
					"maxAttempts": 4,
					"initialBackoff": "0.1s",
					"maxBackoff": "1s",
					"backoffMultiplier": 2.0,
					"retryableStatusCodes": ["UNAVAILABLE"]
				}
			}]
		}`, weightedRoundRobinName)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func main() {
	aAddresses := []Address{
		{Ip: "127.0.0.1", Port: "4002", Weight: 100},
		{Ip: "127.0.0.1", Port: "4003", Weight: 10},
		{Ip: "127.0.0.1", Port: "4004", Weight: 2},
		{Ip: "127.0.0.1", Port: "4005", Weight: 1},
	}

	// IMPORTANT： NewClient 可以一次填入多組位址（例如 4 個 IP:Port + Weight），底層用 manual
	NewClient(aAddresses)
}
