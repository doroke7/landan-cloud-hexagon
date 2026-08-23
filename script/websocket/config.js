'use strict'

// ============================================================
// 跟 config/services.yaml 的 websocket 區塊、config/services.yaml 的
// websocket.admin 是同一組設定；純 demo 用途才直接寫死在這裡，正式環境不會把
// private/public key 跟 signature salt 這樣放進版控。
//
// port 對應 bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT（config/services.yaml
// 的 websocket.port），路徑對應 internal/register/websocket.go 的
// pkg.NewWebsocketRouter("/") ——router 掛在根路徑 "/" 上，不是 "/ws"。
// ============================================================
const WEBSOCKET_URL = 'ws://localhost:4031/'

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
