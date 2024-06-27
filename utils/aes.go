package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Aes struct {
	key string
	iv  string
}

var AesInstant Aes = NewAes()

func generateKey(key, iv []byte) cipher.BlockMode {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return cipher.NewCBCDecrypter(block, iv)
}

func (a Aes) Encrypt(plaintext string) string {
	ivBytes := []byte(a.iv)
	plaintextBytes := []byte(plaintext)

	block, err := aes.NewCipher([]byte(a.key))
	if err != nil {
		panic(err)
	}
	plaintextBytes = pkcs7Pad(plaintextBytes, aes.BlockSize)
	ciphertext := make([]byte, len(plaintextBytes))

	mode := cipher.NewCBCEncrypter(block, ivBytes)
	mode.CryptBlocks(ciphertext, plaintextBytes)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func (a Aes) Decrypt(ciphertext []byte) string {
	ivBytes := []byte(a.iv)
	block, err := aes.NewCipher([]byte(a.key))
	if err != nil {
		panic(err)
	}
	ciphertextBytes := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(ciphertextBytes, ciphertext)

	return string(pkcs5Unpad(ciphertextBytes))
}

func (a Aes) DecryptStr(ciphertext string) string {
	decode, _ := base64.StdEncoding.DecodeString(ciphertext)
	return a.Decrypt(decode)
}
func (a Aes) DecryptBase64(ciphertext []byte) string {
	dst := make([]byte, len(ciphertext))
	base64.StdEncoding.Decode(dst, ciphertext)
	return a.Decrypt(dst)
}

func pkcs7Pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func pkcs5Pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func pkcs5Unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func pkcs7Unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func NewAesByKey(key, iv string) Aes {
	return Aes{
		key: key,
		iv:  iv,
	}
}

func NewAes() Aes {
	return NewAesByKey("ca72ed29dc5eed56b203057f50c6c4de", "0000000000000000")
}
