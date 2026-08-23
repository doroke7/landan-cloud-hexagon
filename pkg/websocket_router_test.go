package pkg

import (
	"testing"

	types "example/types"
)

// TestWebsocketRouterGroupChildrenAndSeparatorMismatch 驗證兩件事：
//  1. 巢狀 Group（children）會依 root -> 祖先 -> 自己 的順序疊加 middleware。
//  2. Group("Admin").HandleFunc("/Authentication/Authenticator.SignIn", ...) 這種
//     混用 "/" 註冊、client 端卻用純 "." 送出 method 的情況，dispatch 要能對到同一個
//     節點——這是把 handler 存進 tree（跟 middleware 共用同一套 tokenize）之後
//     才修好的，改回舊版「routes map[string]handler 用原始字串當 key」就會查不到。
func TestWebsocketRouterGroupChildrenAndSeparatorMismatch(t *testing.T) {
	oRouter := NewWebsocketRouter("/")

	var aOrder []string
	fnMw := func(sTag string) types.WebsocketMiddlewareFunc {
		return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {
			aOrder = append(aOrder, sTag)
			return fnNext(oConn, oReq)
		}
	}

	oRoot := oRouter.Group("", fnMw("root"))
	oAdmin := oRoot.Group("Admin", fnMw("admin"))
	oAuth := oAdmin.Group("Authentication", fnMw("auth"))
	oAuth.HandleFunc("/Authenticator.SignIn", func(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse {
		aOrder = append(aOrder, "handler")
		return types.WebsocketResponse{Code: 1}
	})

	oReq := types.WebsocketRequest{Method: "Admin.Authentication.Authenticator.SignIn"}
	oResp := oRouter.dispatch(nil, oReq)

	if oResp.Code != 1 {
		t.Fatalf("預期 dispatch 找到 handler 並回傳 Code=1，got Code=%d Message=%q", oResp.Code, oResp.Message)
	}

	aExpected := []string{"root", "admin", "auth", "handler"}
	if len(aOrder) != len(aExpected) {
		t.Fatalf("middleware/handler 執行順序長度不對，got %v，want %v", aOrder, aExpected)
	}
	for i, sTag := range aExpected {
		if aOrder[i] != sTag {
			t.Fatalf("middleware/handler 執行順序不對，got %v，want %v", aOrder, aExpected)
		}
	}
}

// TestWebsocketRouterMethodNotFound 驗證沒註冊過的 method 會回傳
// ErrWebsocketMethodNotFound，不會 panic（例如樹狀查找中途 miss 掉子節點）。
func TestWebsocketRouterMethodNotFound(t *testing.T) {
	oRouter := NewWebsocketRouter("/")
	oGroup := oRouter.Group("Admin")
	oGroup.HandleFunc(".Authenticator.SignIn", func(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse {
		return types.WebsocketResponse{Code: 1}
	})

	oResp := oRouter.dispatch(nil, types.WebsocketRequest{Method: "Admin.Authenticator.DoesNotExist"})

	if oResp.Message != ErrWebsocketMethodNotFound.Error() {
		t.Fatalf("預期 method not found 錯誤，got %+v", oResp)
	}
}
