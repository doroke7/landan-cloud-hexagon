'use strict'

const crypto = require('crypto')

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

module.exports = {
  rsaEncrypt,
  aesEncrypt,
  aesDecrypt,
  md5,
  randomString
}
