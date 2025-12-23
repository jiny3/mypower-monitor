package library

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"
)

// 必须与 JS 源码中的 $aes_chars 保持一致
const aesChars = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"

// randomString 还原 JS 中的随机字符串生成逻辑
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = aesChars[rand.Intn(len(aesChars))]
	}
	return string(b)
}

// PKCS7Padding 标准补码
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// encryptPassword 实现完整的加密逻辑
func encryptPassword(password string, salt string) (string, error) {
	// 1. 准备 Key (salt) 和随机生成的 IV
	key := []byte(salt)
	iv := []byte(randomString(16))

	// 2. 构造待加密的明文: 64位随机前缀 + 原始密码
	plainText := []byte(randomString(64) + password)

	// 3. 初始化 AES-CBC 模块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 4. 执行填充并加密
	paddedText := PKCS7Padding(plainText, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)

	encrypted := make([]byte, len(paddedText))
	blockMode.CryptBlocks(encrypted, paddedText)

	// 5. 返回 Base64 编码字符串
	return base64.StdEncoding.EncodeToString(encrypted), nil
}
