package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

type RsaHelper struct {
	*AbstractHelper
}

// NewXYZ 的好處是，init 一個實例的時候，需要哪些基本參數都規定了，不容易遺漏
// 確實 在 New 函數中， 屬性值變成必填了， 但是 直接 new （&）的時候不是

func NewRsaHelper(oAbstractHelper *AbstractHelper) *RsaHelper {
	return &RsaHelper{
		AbstractHelper: oAbstractHelper,
	}
}

// Encrypt encrypts sInput with the provided PEM public key (chunked, 100-byte blocks) and returns base64.
func (oSelf *RsaHelper) Encrypt(sInput, sPublicKey string) (string, error) {
	oPubKey, oErr := parsePublicKey(sPublicKey)
	if oErr != nil {
		return "", fmt.Errorf("rsa encrypt: %w", oErr)
	}

	byInput := []byte(sInput)
	iChunkSize := 100
	var byEncrypted []byte

	for i := 0; i < len(byInput); i += iChunkSize {
		iEnd := i + iChunkSize
		if iEnd > len(byInput) {
			iEnd = len(byInput)
		}

		byChunk, oErr := rsa.EncryptPKCS1v15(rand.Reader, oPubKey, byInput[i:iEnd])
		if oErr != nil {
			return "", fmt.Errorf("rsa encrypt chunk: %w", oErr)
		}

		byEncrypted = append(byEncrypted, byChunk...)
	}

	return base64.RawURLEncoding.EncodeToString(byEncrypted), nil
}

// Decrypt decrypts base64-encoded sInput with the provided PEM private key (chunked by key size).
func (oSelf *RsaHelper) Decrypt(sInput, sPrivateKey string) (string, error) {
	oPrivKey, oErr := parsePrivateKey(sPrivateKey)
	if oErr != nil {
		return "", fmt.Errorf("rsa decrypt: %w", oErr)
	}

	byInput, oErr := base64.RawURLEncoding.DecodeString(strings.TrimRight(sInput, "="))
	if oErr != nil {
		return "", fmt.Errorf("rsa decrypt base64: %w", oErr)
	}

	iChunkSize := oPrivKey.Size()
	var byDecrypted []byte

	for i := 0; i < len(byInput); i += iChunkSize {
		iEnd := i + iChunkSize
		if iEnd > len(byInput) {
			iEnd = len(byInput)
		}

		byChunk, oErr := rsa.DecryptPKCS1v15(rand.Reader, oPrivKey, byInput[i:iEnd])
		if oErr != nil {
			return "", fmt.Errorf("rsa decrypt chunk: %w", oErr)
		}

		byDecrypted = append(byDecrypted, byChunk...)
	}

	return string(byDecrypted), nil
}

func parsePublicKey(sPem string) (*rsa.PublicKey, error) {
	oBlock, _ := pem.Decode([]byte(sPem))
	if oBlock == nil {
		return nil, errors.New("failed to parse PEM block for public key")
	}

	oPub, oErr := x509.ParsePKIXPublicKey(oBlock.Bytes)
	if oErr != nil {
		return nil, fmt.Errorf("parse public key: %w", oErr)
	}

	oRsaPub, bOk := oPub.(*rsa.PublicKey)
	if !bOk {
		return nil, errors.New("key is not RSA public key")
	}

	return oRsaPub, nil
}

func parsePrivateKey(sPem string) (*rsa.PrivateKey, error) {
	oBlock, _ := pem.Decode([]byte(sPem))
	if oBlock == nil {
		return nil, errors.New("failed to parse PEM block for private key")
	}

	oPriv, oErr := x509.ParsePKCS8PrivateKey(oBlock.Bytes)
	if oErr != nil {
		// Fallback to PKCS1
		oPrivPkcs1, oErrPkcs1 := x509.ParsePKCS1PrivateKey(oBlock.Bytes)
		if oErrPkcs1 != nil {
			return nil, fmt.Errorf("parse private key: %w (pkcs1: %w)", oErr, oErrPkcs1)
		}
		return oPrivPkcs1, nil
	}

	oRsaPriv, bOk := oPriv.(*rsa.PrivateKey)
	if !bOk {
		return nil, errors.New("key is not RSA private key")
	}

	return oRsaPriv, nil
}
