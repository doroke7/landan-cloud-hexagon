'use strict'

// 瀏覽器版：Node 的 crypto 模組換成 forge（由 index.html 用 <script> 從 CDN 載入，
// 這裡直接用全域的 window.forge）。演算法、輸出格式跟原本 Node 版完全一致，
// 只是底層實作換了。

// -------------------- 加解密工具，對應 internal/helper/aes_helper.go、rsa_helper.go --------------------

function base64ToBase64Url (sBase64) {
  return sBase64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

function base64UrlToBase64 (sBase64Url) {
  let sBase64 = sBase64Url.replace(/-/g, '+').replace(/_/g, '/')
  while (sBase64.length % 4 !== 0) {
    sBase64 += '='
  }
  return sBase64
}

// 對應 helper.RsaHelper.Encrypt：RSA PKCS1v15，輸出 base64url（不補 = padding）。
function rsaEncrypt (sPlainText, sPublicKeyPem) {
  const oPublicKey = forge.pki.publicKeyFromPem(sPublicKeyPem)
  const sEncrypted = oPublicKey.encrypt(forge.util.encodeUtf8(sPlainText), 'RSAES-PKCS1-V1_5')
  return base64ToBase64Url(forge.util.encode64(sEncrypted))
}

// 對應 helper.AesHelper.Encrypt：AES-128-CBC + PKCS7 padding，輸出 base64url。
// forge 的 cipher 預設就會自動補 PKCS7 padding，不用自己手動補。
function aesEncrypt (sPlainText, sKey, sIv) {
  const oCipher = forge.cipher.createCipher('AES-CBC', forge.util.createBuffer(sKey))
  oCipher.start({ iv: forge.util.createBuffer(sIv) })
  oCipher.update(forge.util.createBuffer(forge.util.encodeUtf8(sPlainText)))
  oCipher.finish()
  return base64ToBase64Url(forge.util.encode64(oCipher.output.getBytes()))
}

// 對應 helper.AesHelper.Decrypt：反過來解密 base64url 密文，空字串直接回傳空字串。
function aesDecrypt (sCipherText, sKey, sIv) {
  if (!sCipherText) {
    return ''
  }
  const sBytes = forge.util.decode64(base64UrlToBase64(sCipherText))
  const oDecipher = forge.cipher.createDecipher('AES-CBC', forge.util.createBuffer(sKey))
  oDecipher.start({ iv: forge.util.createBuffer(sIv) })
  oDecipher.update(forge.util.createBuffer(sBytes))
  oDecipher.finish()
  return forge.util.decodeUtf8(oDecipher.output.getBytes())
}

function md5 (sString) {
  const oMd = forge.md.md5.create()
  oMd.update(forge.util.encodeUtf8(sString))
  return oMd.digest().toHex()
}

function randomString (iLength) {
  // AES-128 的 key/iv 各要 16 bytes；隨便湊可見字元就好，長度對就行。
  // 用 window.crypto.getRandomValues 取真隨機（瀏覽器原生支援，不用額外的庫）。
  const sCharset = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  const aRandomBytes = new Uint8Array(iLength)
  window.crypto.getRandomValues(aRandomBytes)

  let sResult = ''
  for (let i = 0; i < iLength; i++) {
    sResult += sCharset[aRandomBytes[i] % sCharset.length]
  }
  return sResult
}
