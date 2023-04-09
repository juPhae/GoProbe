package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// Encrypt 使用 AES 加密和 CFB 模式以及 Zero Padding 对输入进行加密。
// 它接收一个明文 []byte 和密钥 []byte，并将加密后的密文作为 base64 []byte 和一个错误返回。
func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 从 CFB 模式创建 Cipher Text Stealing (CTS) 模式,即偏移量量是1234567887654321
	iv := []byte("1234567887654321")
	cfb := cipher.NewCFBEncrypter(block, iv)

	// 对输入进行 Zero Padding
	blockSize := block.BlockSize()
	plaintext = zeroPadding(plaintext, blockSize)

	// 加密输入
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)

	// 将密文编码为 base64
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// Decrypt 使用 AES 加密和 CFB 模式以及 Zero Padding 对输入进行解密。
// 它接收一个 base64 []byte 和密钥 []byte，并将解密后的明文作为 []byte 和一个错误返回。

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 从 CFB 模式创建 Cipher Text Stealing (CTS) 模式
	iv := []byte("1234567887654321")
	cfb := cipher.NewCFBDecrypter(block, iv)

	// 将 base64 密文解码为 []byte
	ciphertext, err = base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, err
	}

	// 解密密文
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)

	// 移除 Zero Padding
	plaintext = zeroUnPadding(plaintext)

	return plaintext, nil
}

func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
