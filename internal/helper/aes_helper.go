package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	_ "fmt"
	"strings"
)

/*
*
你用 base64格式，使用postman 的时候，记得把 + / 从 AES-online 换成 - _
gw55ZcBQOW+lgUmjzCRyzA==
gw55ZcBQOW-lgUmjzCRyzA==
  ┌────────────────────┬───────────────────┬────────────────────────────────────┐
  │        编码        │      字符集       │              对应 PHP              │
  ├────────────────────┼───────────────────┼────────────────────────────────────┤
  │ base64.StdEncoding │ A-Z a-z 0-9 + / = │ base64_encode()                    │
  ├────────────────────┼───────────────────┼────────────────────────────────────┤
  │ base64.URLEncoding │ A-Z a-z 0-9 - _ = │ strtr(base64_encode(), '+/', '-_') │
  └────────────────────┴───────────────────┴────────────────────────────────────┘
*/

/*
Aes 128 需要 128 個 bit
然後 一個 英文字 可以 用一個 8bits（1byte） 的 utf8 表示
16 個英文字 剛好 等於 128bits
*/
type AesHelper struct {
	*AbstractHelper
}

func NewAesHelper(oAbstractHelper *AbstractHelper) *AesHelper {
	return &AesHelper{
		AbstractHelper: oAbstractHelper,
	}
}

func (oSelf *AesHelper) Encrypt(sText string, sKey string, sIv string) (string, error) {

	aByteText := []byte(sText)
	aByteKey := []byte(sKey)
	aByteIv := []byte(sIv)

	aByteText = oSelf.pKCS7Padding(aByteText)

	aByteResult := make([]byte, len(aByteText))

	oCipher, oErr := aes.NewCipher(aByteKey)
	if oErr != nil {
		return "", oErr
	}

	oCipherEncrpter := cipher.NewCBCEncrypter(oCipher, aByteIv)
	oCipherEncrpter.CryptBlocks(aByteResult, aByteText)

	sResult := base64.RawURLEncoding.EncodeToString(aByteResult)

	return sResult, nil
}

func (oSelf *AesHelper) Decrypt(sText string, sKey string, sIv string) (string, error) {

	if sText == "" {
		return "", nil
	}

	// IMPORTANT: 前端進來的 加密 base64 ， 先 ”統一“ 把它強制轉換 成 url-base64， 再動作
	sText = strings.TrimRight(sText, "=")
	oReplacer := strings.NewReplacer("+", "-", "/", "_")
	sText = oReplacer.Replace(sText)
	sUrlBase64 := sText

	aByteText, oDecodeErr := base64.RawURLEncoding.DecodeString(sUrlBase64)
	if oDecodeErr != nil {
		return "", oDecodeErr
	}

	aByteKey := []byte(sKey)
	aByteIv := []byte(sIv)

	oBlock, oErr := aes.NewCipher(aByteKey)
	if oErr != nil {
		return "", oErr
	}

	aByteResult := make([]byte, len(aByteText))

	oCipher := cipher.NewCBCDecrypter(oBlock, aByteIv)
	oCipher.CryptBlocks(aByteResult, aByteText)

	aByteResult = oSelf.pKCS7UnPadding(aByteResult)

	sResult := string(aByteResult)

	return sResult, nil
}

func (oSelf *AesHelper) pKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	aPadding := []byte{byte(padding)}
	padtext := bytes.Repeat(aPadding, padding)
	aResult := append(ciphertext, padtext...)

	return aResult
}

func (oSelf *AesHelper) pKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
