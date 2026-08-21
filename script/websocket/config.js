'use strict'

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

module.exports = {
  WEBSOCKET_URL,
  RSA_PUBLIC_KEY,
  SIGNATURE_SALT,
  TYPE_EVENT,
  TYPE_ACK
}
