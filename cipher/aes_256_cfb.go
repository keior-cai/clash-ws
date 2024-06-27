package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type Aes256CfbCipher struct {
	password     string
	decodeStream cipher.Stream
	encodeStream cipher.Stream
	keyLength    int
}

func NewAes256CfbCipher(password string) *Aes256CfbCipher {
	return &Aes256CfbCipher{
		password:  password,
		keyLength: 32,
	}
}

func (c *Aes256CfbCipher) Encrypt(plaintext []byte) ([]byte, error) {
	var target []byte
	offset := 0
	if c.encodeStream == nil {
		block, err := aes.NewCipher(getKey(c.password, c.keyLength))
		if err != nil {
			return nil, err
		}
		iv := make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
		offset = aes.BlockSize
		target = make([]byte, len(plaintext)+aes.BlockSize)
		copy(target, iv)
		c.encodeStream = cipher.NewCFBEncrypter(block, iv)
	} else {
		target = make([]byte, len(plaintext))
	}
	c.encodeStream.XORKeyStream(target[offset:], plaintext)
	return target, nil
}

func (c *Aes256CfbCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	var target []byte
	offset := 0
	if c.decodeStream == nil {
		block, err := aes.NewCipher(getKey(c.password, c.keyLength))
		if err != nil {
			return nil, err
		}
		if len(ciphertext) < aes.BlockSize {
			return nil, fmt.Errorf("ciphertext too short")
		}
		iv := ciphertext[:aes.BlockSize]
		offset = aes.BlockSize
		c.decodeStream = cipher.NewCFBDecrypter(block, iv)
	}
	target = make([]byte, len(ciphertext)-offset)
	c.decodeStream.XORKeyStream(target, ciphertext[offset:])
	return target, nil
}
