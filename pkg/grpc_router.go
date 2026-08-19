package pkg

import (
	"context"
	"strings"

	"google.golang.org/grpc"
)

/*
GrpcRouter 模仿 gin 巢狀 group 的效果：prefix 用「.」「/」分段建成一棵樹，
dispatch 時沿著 info.FullMethod 的每一段往下走，沿路每經過一個真的被
Group() 註冊過的節點，就把該節點自己的攔截器疊加上去——跟 gin 的
parentGroup.Group(child, middlewareB) 一樣，子路徑會自動繼承所有祖先
group 的攔截器，再疊加自己的，順序是「祖先在前、自己在後」，
且跟 Group() 呼叫的先後順序完全無關。

沒有獨立的「全局 aBase」概念：root 本身就是樹的一個節點，想要「不管打
哪個 method 都要跑」的全局攔截器，直接 Group("", 全局攔截器...) 註冊在
root 上即可，dispatch 時一定會先經過 root，效果跟以前的 aBase 一樣，
只是統一成同一套機制，不用額外的欄位。
*/

type routeNode struct {
	children     map[string]*routeNode
	interceptors []grpc.UnaryServerInterceptor // 這個節點自己註冊的攔截器，不含 aBase、也不含祖先節點的
	registered   bool                          // 這個節點是不是真的被 Group() 註冊過，還是只是路過的中繼節點
}

func newRouteNode() *routeNode {
	return &routeNode{children: make(map[string]*routeNode)}
}

type GrpcRouter struct {
	root *routeNode
}

func NewGrpcRouter() *GrpcRouter {
	return &GrpcRouter{root: newRouteNode()}
}

func ChainInterceptors(aInterceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(oContext context.Context, oRequest any, oServerInfo *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (any, error) {

		fnNextHandler := fnHandler

		for i := len(aInterceptors) - 1; i >= 0; i-- {
			i, fnThisHandler := i, fnNextHandler
			fnNextHandler = func(oCtx context.Context, oReq any) (any, error) {
				return aInterceptors[i](oCtx, oReq, oServerInfo, fnThisHandler)
			}
		}

		return fnNextHandler(oContext, oRequest)
	}
}
func tokenize(sPath string) []string {
	return strings.FieldsFunc(sPath, func(r rune) bool {
		return r == '.' || r == '/'
	})
}

// Group 每次 aaa.bbb.ccc 的 path 就生成 tree 的結構，並且在最後的 葉子 加上 middlewares
func (oSelf *GrpcRouter) Group(sPrefix string, aInterceptors ...grpc.UnaryServerInterceptor) *GrpcRouter {
	oNode := oSelf.root
	for _, sSegment := range tokenize(sPrefix) {
		oChild, bOk := oNode.children[sSegment]
		if !bOk {
			oChild = newRouteNode()
			oNode.children[sSegment] = oChild
		}
		oNode = oChild
	}
	oNode.interceptors = aInterceptors
	oNode.registered = true

	return oSelf
}

// 根據 aaa.bbb.ccc 查找 每個子 path 的 middleware 並且合併
func (oSelf *GrpcRouter) Build() grpc.UnaryServerInterceptor {
	return func(oContex context.Context, oRequest any, oServerInfo *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (any, error) {
		var aInterceptors []grpc.UnaryServerInterceptor
		oNode := oSelf.root
		if oNode.registered {
			aInterceptors = append(aInterceptors, oNode.interceptors...)
		}

		for _, sSegment := range tokenize(oServerInfo.FullMethod) {
			oChild, bOk := oNode.children[sSegment]
			if !bOk {
				break
			}
			oNode = oChild
			if oNode.registered {
				aInterceptors = append(aInterceptors, oNode.interceptors...)
			}
		}

		return ChainInterceptors(aInterceptors...)(oContex, oRequest, oServerInfo, fnHandler)
	}
}
