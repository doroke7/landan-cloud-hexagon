'use strict'

const { WEBSOCKET_URL } = require('./config')
const SocketClient = require('./socket-client')
const SocketRouter = require('./socket-router')

// -------------------- 主流程 --------------------

function main () {
  const sName = process.argv[2] || 'admin'
  const sPassword = process.argv[3] || 'password'

  const oSocketRouter = new SocketRouter()
  const oClient = new SocketClient(WEBSOCKET_URL, oSocketRouter)

  oClient.onRoute('Admin.Authentication.Authenticator.SignIn.Welcome', (oParam) => {
    console.log('--- 收到 Welcome event ---')
    console.log(oParam)
  })

  oClient.ready(() => {
    console.log(`已連線到 ${WEBSOCKET_URL}`)

    oClient.AdminAuthenticationAuthenticatorSignIn({ name: sName, password: sPassword }, (oResult) => {
      console.log('--- SignIn 回應 ---')
      console.log(oResult)
    })
  })
}

main()
