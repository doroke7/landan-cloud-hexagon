package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

// PKCS7Padding 填充至 16 字节倍数
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 去除填充
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("数据为空")
	}
	unpadding := int(origData[length-1])
	if unpadding > length {
		return nil, errors.New("解密填充格式错误")
	}
	return origData[:(length - unpadding)], nil
}

// EncryptCBC AES-CBC 加密 (显式传入 IV)
func EncryptCBC(plaintext string, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	// 校验 IV 长度
	if len(iv) != blockSize {
		return "", fmt.Errorf("IV 长度必须为 %d 字节", blockSize)
	}

	// 1. 进行填充
	content := PKCS7Padding([]byte(plaintext), blockSize)

	// 2. 准备密文空间
	ciphertext := make([]byte, len(content))

	// 3. 开始加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, content)

	// 使用 URL 安全编码返回
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptCBC AES-CBC 解密 (显式传入 IV)
func DecryptCBC(cryptoText string, key []byte, iv []byte) (string, error) {
	decodeData, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return "", fmt.Errorf("IV 长度必须为 %d 字节", blockSize)
	}

	// 1. 直接解密整个数据 (因为 IV 是外部传入的，密文里不再包含 IV)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decodeData, decodeData)

	// 2. 去除填充
	result, err := PKCS7UnPadding(decodeData)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func main() {
	// 1. 密钥：16, 24, 32 字节
	key := []byte("12345678901234567890123456789012")

	// 2. IV：固定 16 字节 (配合你的 dKeys 结构体使用)
	iv := []byte("1234567890123456")

	originalText := "Gin Backend Architect 2026"

	// 加密
	encryptStr, err := EncryptCBC(originalText, key, iv)
	if err != nil {
		fmt.Println("加密失败:", err)
		return
	}
	fmt.Println("CBC 加密密文:", encryptStr)

	// 解密
	decryptStr, err := DecryptCBC(encryptStr, key, iv)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}
	fmt.Println("CBC 解密结果:", decryptStr)
}
