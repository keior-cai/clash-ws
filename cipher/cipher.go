package cipher

import "crypto/md5"

type Cipher interface {
	Decrypt([]byte) ([]byte, error)

	Encrypt([]byte) ([]byte, error)
}

func getKey(password string, keyLength int) []byte {
	result := make([]byte, keyLength)

	// 创建MD5哈希对象
	messageDigest := md5.New()

	for hasLength := 0; hasLength < keyLength; hasLength += 16 {
		passwordBytes := []byte(password)

		// 组合需要摘要的字节
		combineBytes := make([]byte, hasLength+len(passwordBytes))
		copy(combineBytes, result[:hasLength])
		copy(combineBytes[hasLength:], passwordBytes)

		// 哈希组合字节
		messageDigest.Reset()
		messageDigest.Write(combineBytes)
		digestBytes := messageDigest.Sum(nil)

		addLength := 16
		if hasLength+16 > keyLength {
			addLength = keyLength - hasLength
		}
		copy(result[hasLength:], digestBytes[:addLength])
	}

	return result
}

func NewCipher(name, password string) Cipher {
	switch name {
	case "aes-256-cfb":
		return NewAes256CfbCipher(password)
	case "aes-128-cfb":
		return NewAes128CfbCipher(password)
	}
	return nil
}
