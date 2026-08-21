'use strict'

const { encodeRequest, decryptResponse } = require('./protocol')

// -------------------- SocketClient --------------------

// 一般
// client 對 server 大部分 聊天訊息 需要 ack
// client 對 server 小部分 打字狀態 不要 ack
// server 對 client 大部分 消息。   不要 ack

class SocketClient {
  constructor (sUrl, oSocketRouter, oSocketController) {

    this.socket = new WebSocket(sUrl)
    this.requestId = 0
    this.callbacks = new Map()
    this.socketRouter = oSocketRouter
    this.socketController = oSocketController

    this.socket.addEventListener('message', (oEvent) => {
      const oResponse = JSON.parse(oEvent.data)

      if (oResponse.type == 'event') {
        this.socketRouter.dispatch(oResponse)
        return
      }

      const iRequestId = oResponse['request-id']
      const oCallback = this.callbacks.get(iRequestId)
      if (!oCallback) {
        return
      }
      this.callbacks.delete(iRequestId)
      // 只要 server 有回任何一種型別的回應（ack 或 normal），就代表這筆
      // request 已經送達、不算「沒收到回應」，停止重送計時器。
      clearInterval(oCallback.timer)

      if(oResponse.type == 'ack') {
        oCallback.callback(decryptResponse(oResponse, oCallback.key, oCallback.iv))
        return
      }

      if(oResponse.type == 'normal') {
        // DO NOTHING
        // normal 不 callback
        console.error('normal 不 callback')

        return
      }


    })

    this.socket.addEventListener('error', (oEvent) => {
      console.error('websocket error:', oEvent.error || oEvent.message)
    })
  }

  // 連線建立完成才能送第一筆訊息，呼叫端要在 fnCallback 裡面才開始呼叫其他方法。
  ready (fnCallback) {
    this.socket.addEventListener('open', () => fnCallback(), { once: true })
  }

  close () {
    this.socket.close()
  }

  // 沒收到 server 回應（callback 一直沒被觸發）就每隔 iRetryInterval ms 重送同一筆
  // 訊息（同一個 request-id），直到收到回應為止——用來扛連線抖動、訊息丟失這類情況，
  // 不是無限期等待。收到回應（不管是 ack 還是 normal）就會在 message listener 裡
  // clearInterval，停止重送。
  emit (sMethod, oParam, fnCallback, iRetryInterval = 2000) {
    const iRequestId = ++this.requestId
    const { key, iv, message } = encodeRequest(iRequestId, sMethod, oParam)
    const sPayload = JSON.stringify(message)

    const fnSend = () => this.socket.send(sPayload)

    if (fnCallback) {
      const oTimer = setInterval(fnSend, iRetryInterval)
      this.callbacks.set(iRequestId, { key, iv, callback: fnCallback, timer: oTimer })
    }

    fnSend()
  }


  AdminAuthenticationAuthenticatorSignIn (oParam, fnCallback) {
    this.emit('Admin.Authentication.Authenticator.SignIn', oParam, fnCallback)
  }
}

module.exports = SocketClient
