'use strict'

/**
   Websocket 的問題



 * 
 * 
 */

const crypto = require('crypto')

// ============================================================
// 跟 config/websocket.yaml、config/admin.yaml 是同一組設定；
// 純 demo 用途才直接寫死在這裡，正式環境不會把 private/public key
// 跟 signature salt 這樣放進版控。
// ============================================================
const WEBSOCKET_URL = 'ws://localhost:8989/ws'

const RSA_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqBbztPHxvEsnb5BcMTjv
693XqRMYte+ORJvrgc0RsdHlC4W4lSqvjnH2JcsTGgtSqmmvsvoPQ9dKGs6+OD3E
zIClDiEr4n7QRFFjYKP3IBqhkR5a5wZdiOCoYCx2dOKjBTkLgIzMO145nITHR0za
Yv7k22eNdIzlLVat1Oq1DlWCHWBEQHUUm/OhiBSHnRb2DXiMa+vBvHHrZBIcDb0+
TRD14zLArY5ijKWkzTLGzr4IDi3TcwDz6xEkLm4grzi/KEYtjAweVTClqm19vYAk
SDe+BtVYNxODv3yQSSIrDEzeCnbimIBCBfwxL65YrbIAUx7YqVbtNry56C4MI95h
rQIDAQAB
-----END PUBLIC KEY-----`

const SIGNATURE_SALT = '~9U7g2R8zgW&dZ_u'

// 對應 types.WebsocketRequest/Response.Type：client 送出的請求固定是 "event"，
// server 回應固定是 "ack"（見 pkg.WebsocketRouter.serveConn），跟
// sample/websocket/client.js 的 WSClient 用法一致。
const TYPE_EVENT = 'event'
const TYPE_ACK = 'ack'

// -------------------- 加解密工具，對應 internal/helper/aes_helper.go、rsa_helper.go --------------------

// 對應 helper.RsaHelper.Encrypt：RSA PKCS1v15，輸出 base64url（不補 = padding）。
function rsaEncrypt (sPlainText, sPublicKeyPem) {
  const oEncrypted = crypto.publicEncrypt(
    { key: sPublicKeyPem, padding: crypto.constants.RSA_PKCS1_PADDING },
    Buffer.from(sPlainText, 'utf8')
  )
  return oEncrypted.toString('base64url')
}

// 對應 helper.AesHelper.Encrypt：AES-128-CBC + PKCS7 padding，輸出 base64url。
// createCipheriv 預設就會自動補 PKCS7 padding，不用自己手動補。
function aesEncrypt (sPlainText, sKey, sIv) {
  const oCipher = crypto.createCipheriv('aes-128-cbc', Buffer.from(sKey, 'utf8'), Buffer.from(sIv, 'utf8'))
  const oEncrypted = Buffer.concat([oCipher.update(sPlainText, 'utf8'), oCipher.final()])
  return oEncrypted.toString('base64url')
}

// 對應 helper.AesHelper.Decrypt：反過來解密 base64url 密文，空字串直接回傳空字串。
function aesDecrypt (sCipherText, sKey, sIv) {
  if (!sCipherText) {
    return ''
  }
  const oDecipher = crypto.createDecipheriv('aes-128-cbc', Buffer.from(sKey, 'utf8'), Buffer.from(sIv, 'utf8'))
  const oDecrypted = Buffer.concat([oDecipher.update(Buffer.from(sCipherText, 'base64url')), oDecipher.final()])
  return oDecrypted.toString('utf8')
}

function md5 (sString) {
  return crypto.createHash('md5').update(sString, 'utf8').digest('hex')
}

function randomString (iLength) {
  // AES-128 的 key/iv 各要 16 bytes；隨便湊可見字元就好，長度對就行。
  const sCharset = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let sResult = ''
  for (let i = 0; i < iLength; i++) {
    sResult += sCharset[crypto.randomInt(sCharset.length)]
  }
  return sResult
}

// 對應 internal/middleware/websocket/admin 那條 Signature -> Decryption -> ... -> Encryption
// chain：這裡在 client 端做「反過來」的事——加密 param、算簽名，讓 server 那條 chain
// 可以正確驗簽、解密。
//
// iRequestId 對應 types.WebsocketRequest.RequestId：server 的 pkg.WebsocketRouter.serveConn
// 會原樣把它放進回傳的 WebsocketResponse，client 端就不用假設「先送先回」，直接用
// requestId 把回應配對回正確的呼叫方（跟 sample/websocket/client.js 的做法一致）。
function encodeRequest (iRequestId, sMethod, oParam) {
  const sKey = randomString(16)
  const sIv = randomString(16)

  const sK = rsaEncrypt(JSON.stringify({ key: sKey, iv: sIv }), RSA_PUBLIC_KEY)
  const sP = aesEncrypt(JSON.stringify(oParam), sKey, sIv)
  const sS = ''
  const sO = ''
  const sVer = '1.0'
  const sVersion = '1.0'
  const sTime = Date.now().toString()

  // 跟 internal/middleware/websocket/admin/signature_middleware.go 的
  // aStrings := []string{Ver, Version, K, Time, S, O, P, SALT} 順序要完全一致，
  // 不然 SignatureMiddleware 那邊算出來的 md5 對不上。
  const sSignature = md5([sVer, sVersion, sK, sTime, sS, sO, sP, SIGNATURE_SALT].join('|'))

  return {
    key: sKey,
    iv: sIv,
    message: {
      'request-id': iRequestId,
      type: TYPE_EVENT,
      header: { ver: sVer, version: sVersion, k: sK, a: '', time: sTime, signature: sSignature },
      code: 0,
      method: sMethod,
      p: sP,
      s: sS,
      o: sO
    }
  }
}

// 對應 types.WebsocketResponse：DEBUG 模式下（config/default.yaml 的 debug: true）
// code/message/result 會保留明文方便直接看；c/m/r 才是真正的密文，用發起這次呼叫時
// 產生的 key/iv 解開，驗證整個加解密流程真的有跑通。
function decryptResponse (oResp, sKey, sIv) {
  return {
    ...oResp,
    decrypted: {
      code: aesDecrypt(oResp.c, sKey, sIv),
      message: aesDecrypt(oResp.m, sKey, sIv),
      result: aesDecrypt(oResp.r, sKey, sIv)
    }
  }
}

// -------------------- EventRouter --------------------

// EventRouter 對應 server 端 pkg.WebsocketRouter：用 method 當 key 找對應的處理方法，
// 只是這裡分派的是 server 主動推播的 "event"（例如 Admin.Authentication.Authenticator.
// SignIn.Welcome），不是呼叫的回應（ack，走 SocketClient.callbacks 那條路）。
class EventRouter {
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

// -------------------- SocketClient --------------------

// 一般
// client 對 server 大部分 聊天訊息 需要 ack
// client 對 server 小部分 打字狀態 不要 ack
// server 對 client 大部分 消息。   不要 ack

class SocketClient {
  constructor (sUrl) {
    this.socket = new WebSocket(sUrl)
    this.requestId = 0
    this.callbacks = new Map()
    this.eventRouter = new EventRouter()

    this.socket.addEventListener('message', (oEvent) => {
      const oResp = JSON.parse(oEvent.data)

      if (oResp.type == 'event') {
        this.eventRouter.dispatch(oResp)
        return
      }

      const iRequestId = oResp['request-id']
      const oCallback = this.callbacks.get(iRequestId)
      if (!oCallback) {
        return
      }
      this.callbacks.delete(iRequestId)
      // 只要 server 有回任何一種型別的回應（ack 或 normal），就代表這筆
      // request 已經送達、不算「沒收到回應」，停止重送計時器。
      clearInterval(oCallback.timer)

      if(oResp.type == 'ack') {
        oCallback.callback(decryptResponse(oResp, oCallback.key, oCallback.iv))
        return
      }

      if(oResp.type == 'normal') {
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

  // 對應 server 端 pkg.WebsocketRouter.HandleFunc：註冊一個 method 收到 server 主動
  // 推播的 "event" 時要怎麼處理，跟 emit() 是相反方向（emit 是送出去等 ack，onEvent
  // 是被動等 server 推進來）。
  onRoute (sMethod, fnHandler) {
    this.eventRouter.handle(sMethod, fnHandler)
  }

  // 對應 server 端 internal/register/websocket.go 的
  // oRouter.Group("Admin.").HandleFunc("Authentication.Authenticator.SignIn", ...)
  AdminAuthenticationAuthenticatorSignIn (oParam, fnCallback) {
    this.emit('Admin.Authentication.Authenticator.SignIn', oParam, fnCallback)
  }
}

// -------------------- 主流程 --------------------

function main () {
  const sName = process.argv[2] || 'admin'
  const sPassword = process.argv[3] || 'password'

  const oClient = new SocketClient(WEBSOCKET_URL)

  // 對應 internal/input/application/websocket/admin/authentication/authenticator_handler.go
  // 的 SignIn 在登入成功後 oConn.Push("Admin.Authentication.Authenticator.SignIn.Welcome", ...)。
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
