package utils

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	testContent := "123456"
	fmt.Println("加密前:", testContent)
	defaultKey := "ca72ed29dc5eed56b203057f50c6c4de"
	iv := "0000000000000000"

	aes := NewAesByKey(defaultKey, iv)

	content := aes.Encrypt(testContent)
	fmt.Println("加密后的数据", content)
	content = aes.DecryptStr(content)
	fmt.Println("解密后的数据", content)
	content = aes.DecryptStr("YfTwsydT2S7E6qXipB9itz3HhHx3aGkstvypabnvU9Ozfsd67p1RhlSKjUC/BmcL/EoFJo75js3KbpWe/3ud4aUi5GIFcbTx0tcYfTgt3QZxntCnyh2BiakWksH8f+F2MLLCNv5AkTSuSz/k/HTuxDWriXi2NxkA3xYnqfC4WklP8pILO6AOxS4lm+E7p/17en89V1vAIjaRHqKVb21RcJT8YyQb5pho+jZ0PvFfb5/6qCJ2APtGPHvcyLhh+e6UDNzfP8Bo34GplVYtgFtGhnHc3FzQZ7qaxZQIR/Fb7phGDO9xSkbvKR6ML1av0gVavXHoHxPa+1Hum6lwjiQFcDJ519jij8DKRZdZpHBQB1I=")
	fmt.Println("解密后的数据", content)

	//publicKey := encryptByDefaultKey(testContent defaultKey iv)
	//fmt.Println("publicKey 加密后:" publicKey)

	// 这里假设你已经从 Java 中获取了加密后的字符串
	// 如果你有加密后的字符串，你可以直接传递给 decrypt 函数
	// 如果你想手动截取 Java 加密结果的 "="，你也可以这么做
	// 注意：Java 和 Go 在 Base64 编码上可能存在一些差异，需要注意处理
	//publicKey = "Ua25HPj/9u1gnDIMB9PdcELLjycWaNl4VXE610trMU8"
	//
	//privateKey := decrypt(publicKey defaultKey iv)
	//fmt.Println("privatekey 解密后:" content)
}
