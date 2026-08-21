'use strict'

const { WEBSOCKET_URL } = require('./config')
const SocketClient = require('./socket-client')
const SocketRouter = require('./socket-router')
const SocketController = require('./socket-controller')

// -------------------- 主流程 --------------------

function main () {
  const sName = process.argv[2] || 'admin'
  const sPassword = process.argv[3] || 'password'

  const oSocketRouter = new SocketRouter()
  const oSocketController = new SocketController()
  const oClient = new SocketClient(WEBSOCKET_URL, oSocketRouter, oSocketController)

  oSocketRouter.handle('Admin.Authentication.Authenticator.SignIn.Welcome', oSocketController.adminAuthenticationAuthenticatorSignInWelcome)

  oClient.ready(() => {
    console.log(`已連線到 ${WEBSOCKET_URL}`)

    oClient.AdminAuthenticationAuthenticatorSignIn({ name: sName, password: sPassword }, (oResult) => {
      console.log('--- SignIn 回應 ---')
      console.log(oResult)
    })
  })
}

main()
