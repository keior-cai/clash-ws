package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

type RsaUtil struct {
	privateKey *rsa.PrivateKey

	publicKey *rsa.PublicKey

	PublicKey string
}

var RsaInstant = NewRsa("MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDlwS6f4FBSHKDgg8Tti2YXW6ic8BGLeoKI8IuXEUy0q2cV53DcJ7ON55oXuuDuBRLE6PanT86gcoRTp1IOTKjI7fga3arIaWjYubEBzCLUlTPQx/jjO0/mWarj4yvKzk6Ulo/uXWumR+dx0dYiGtbJQlClgILvYtxNHQB7uXWPjwIDAQAB", "MIICeQIBADANBgkqhkiG9w0BAQEFAASCAmMwggJfAgEAAoGBAOXBLp/gUFIcoOCDxO2LZhdbqJzwEYt6gojwi5cRTLSrZxXncNwns43nmhe64O4FEsTo9qdPzqByhFOnUg5MqMjt+BrdqshpaNi5sQHMItSVM9DH+OM7T+ZZquPjK8rOTpSWj+5da6ZH53HR1iIa1slCUKWAgu9i3E0dAHu5dY+PAgMBAAECgYEAk87uYeh7g/fq/8WGAZR2v3w2Q5CmmObd559pDm0QvgKvNQZKMzhPaXGgTrfpUPdulcOSOx06vzotK2wvfAeRZUqmApZqlLOiNkcrafEIBjwBlWh7EKxw9bXauKgdQXr7MPfQg11ipbw52wGXmEElvB5tEuCX5tVD9KHzkluXyKECQQD3fq/WgaDTWlnTVb1QvnyBP+bS+40A9JMst9WDK1qKQ8urKFX4Lnfw7s5953Lbx/euLzM1+e9tnWmcUTMa0Op5AkEA7aZtUg48z9rTN4OposMITmOaO870CZot8DE0RS1MshVsSCL6AbKRFiOLzxoDlFGEFtAvephN5qHPtYhWT4bIRwJBAMbY25Ad4EhPjGIWvh9UnJX/8IXNJBIDbwf7v6k+uOTj6YxfwQrA0w8Z34Aa6BabSG2DcMLKR8srMQIt30CJYAkCQQDI6cDWdG75Evkqn8cUcWpeS1qjYa1zSMO5ov+b1FZY4D+xJNDUCpEadGbIaifIhrnzR4I8VPLXHsmpoV/G0B4VAkEArhaCTjjg5KIyyccBIcyTo8RVCQV1/cEwtdl/b+E4JzFatkMvVLbWVSZJ+b0ZxRqDA4DD6qFaZKl2Ya0vgtiPhg==")

func (r RsaUtil) RsaEncrypt(data []byte) string {
	v15, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, data)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(v15)
}

func (r RsaUtil) RsaEncryptKey(data []byte, publicKey string) string {
	v15, err := rsa.EncryptPKCS1v15(rand.Reader,
		generatePublicKey(publicKey),
		data)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(v15)
}

func (r RsaUtil) RsaEncryptStr(data string) string {
	return r.RsaEncrypt([]byte(data))
}

func (r RsaUtil) RsaDecrypt(encryptedData []byte) (string, error) {
	v15, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, encryptedData)
	return string(v15), err
}

func (r RsaUtil) RsaDecryptStr(encryptedData string) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	return r.RsaDecrypt(decodeString)
}

// 生成公钥
func generatePublicKey(key string) *rsa.PublicKey {

	// 解码 Base64 编码的字符串
	bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return nil
	}
	pubKey, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		fmt.Println("Error parsing public key:", err)
		return nil
	}
	rsaPrivateKey, _ := pubKey.(*rsa.PublicKey)
	return rsaPrivateKey
}

// 生成密钥
func generatePrivateKey(key string) *rsa.PrivateKey {

	// 解码 Base64 编码的字符串
	bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return nil
	}
	// 解析私钥
	privateKey, err := x509.ParsePKCS8PrivateKey(bytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return nil
	}
	rsaPrivateKey, _ := privateKey.(*rsa.PrivateKey)
	return rsaPrivateKey
}

func NewRsa(public, private string) RsaUtil {
	publicKey := generatePublicKey(public)
	privateKey := generatePrivateKey(private)
	return RsaUtil{
		publicKey:  publicKey,
		privateKey: privateKey,
		PublicKey:  public,
	}
}
