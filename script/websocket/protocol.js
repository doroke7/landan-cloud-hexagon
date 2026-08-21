'use strict'

const { RSA_PUBLIC_KEY, SIGNATURE_SALT, TYPE_EVENT } = require('./config')
const { rsaEncrypt, aesEncrypt, aesDecrypt, md5, randomString } = require('./crypto')

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

module.exports = {
  encodeRequest,
  decryptResponse
}
