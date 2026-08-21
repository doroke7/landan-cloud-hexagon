'use strict'

// -------------------- SocketController --------------------

class SocketController {
  constructor () {
    this.adminAuthenticationAuthenticatorSignInWelcome = this.adminAuthenticationAuthenticatorSignInWelcome.bind(this)
  }

  adminAuthenticationAuthenticatorSignInWelcome(oParam) {
    console.log('--- 收到 Welcome event ---')
    console.log(oParam)
  }
}

module.exports = SocketController
