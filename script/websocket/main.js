'use strict'

// -------------------- 主流程（瀏覽器版） --------------------
// 原本 Node 版直接在 process.argv 拿 name/password、跑完就結束；
// 瀏覽器版换成「連線」「SignIn」兩個按鈕各自觸發，結果印到頁面的 log 區塊
// 而不是 console.log，其他協定/加解密邏輯完全沒變。

let oContainer = null

function appendLog (sType, sTitle, oData) {
  const oLog = document.getElementById('log')
  const oLine = document.createElement('div')
  oLine.className = 'log-line log-' + sType

  const sTime = new Date().toLocaleTimeString()
  let sText = `[${sTime}] ${sTitle}`
  if (oData !== undefined) {
    sText += ' ' + JSON.stringify(oData, null, 2)
  }
  oLine.textContent = sText

  oLog.appendChild(oLine)
  oLog.scrollTop = oLog.scrollHeight
}

function setStatus (sText) {
  document.getElementById('status').textContent = sText
}

function connect () {
  if (oContainer) {
    appendLog('info', '已經連線過了，不用重複連線')
    return
  }

  oContainer = initContainer()

  // 對應 main.js（Node 版）的 container.router.handle(...)：註冊 server 主動
  // 推播的 Welcome event。
  oContainer.router.handle(
    'Admin.Authentication.Authenticator.SignIn.Welcome',
    oContainer.controller.adminAuthenticationAuthenticatorSignInWelcome
  )

  setStatus('連線中...')
  appendLog('info', `連線到 ${WEBSOCKET_URL}`)

  oContainer.client.ready(() => {
    setStatus('已連線')
    appendLog('info', '已連線')
    document.getElementById('signin-btn').disabled = false
  })
}

function signIn () {
  if (!oContainer) {
    appendLog('error', '尚未連線，請先按「連線」')
    return
  }

  const sName = document.getElementById('name').value || 'admin'
  const sPassword = document.getElementById('password').value || 'password'

  appendLog('request', 'SignIn 送出', { name: sName })

  oContainer.client.AdminAuthenticationAuthenticatorSignIn({ name: sName, password: sPassword }, (oResult) => {
    appendLog('response', 'SignIn 回應', oResult)
  })
}

window.addEventListener('DOMContentLoaded', () => {
  document.getElementById('connect-btn').addEventListener('click', connect)
  document.getElementById('signin-btn').addEventListener('click', signIn)
})
