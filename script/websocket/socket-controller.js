'use strict'

// -------------------- SocketController --------------------
// server 主動推播的 event 就是在這裡處理；appendLog 是 main.js 提供的 DOM
// 輸出函式，這裡只管收到什麼、要顯示什麼，不管怎麼畫。

class SocketController {
  constructor () {
    this.adminAuthenticationAuthenticatorSignInWelcome = this.adminAuthenticationAuthenticatorSignInWelcome.bind(this)
  }

  adminAuthenticationAuthenticatorSignInWelcome (oParam) {
    appendLog('event', '收到 Welcome event', oParam)
  }
}
