'use strict'

const { WEBSOCKET_URL } = require('./config')
const initContainer = require('./container')

// -------------------- 主流程 --------------------

function main () {
  const sName = process.argv[2] || 'admin'
  const sPassword = process.argv[3] || 'password'

  const container = initContainer();

  container.client.ready(() => {
    console.log(`已連線到 ${WEBSOCKET_URL}`)

    container.client.AdminAuthenticationAuthenticatorSignIn({ name: sName, password: sPassword }, (oResult) => {
      console.log('--- SignIn 回應 ---')
      console.log(oResult)
    })
  })
}

main()
