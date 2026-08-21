'use strict'

// SocketRouter 對應 server 端 pkg.WebsocketRouter：用 method 當 key 找對應的處理方法，
// 只是這裡分派的是 server 主動推播的 "event"（例如 Admin.Authentication.Authenticator.
// SignIn.Welcome），不是呼叫的回應（ack，走 SocketClient.callbacks 那條路）。
class SocketRouter {
  constructor () {
    this.routes = new Map()
  }

  // 對應 pkg.WebsocketRouter.HandleFunc。
  handle (sMethod, fnHandler) {
    this.routes.set(sMethod, fnHandler)
  }

  // 對應 pkg.WebsocketRouter.dispatch：oPacket 是一筆 { method, param } 事件。
  dispatch (oPacket) {
    const fnHandler = this.routes.get(oPacket.method)
    if (!fnHandler) {
      console.log('unhandled event:', oPacket)
      return
    }
    fnHandler(oPacket.param)
  }
}

module.exports = SocketRouter
