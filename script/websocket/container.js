'use strict'

const { WEBSOCKET_URL } = require('./config')
const SocketClient = require('./socket-client')
const SocketRouter = require('./socket-router')
const SocketController = require('./socket-controller')

// -------------------- 依賴注入容器 --------------------

// 組裝 SocketRouter -> SocketController -> SocketClient 的注入鏈，並把 route
// 註冊也收在這裡；main.js 只需要拿組好的 oClient 出來跑業務流程，不用知道
// 這些物件是怎麼兜起來的。
function initContainer () {
  const oSocketRouter = new SocketRouter()
  const oSocketController = new SocketController()
  const oSocketClient = new SocketClient(WEBSOCKET_URL, oSocketRouter, oSocketController)

  let oContainer = {
    router: oSocketRouter, 
    controller: oSocketController, 
    client: oSocketClient
  };


  return oContainer
}

module.exports = initContainer
