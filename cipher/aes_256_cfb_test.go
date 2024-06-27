package cipher

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAes256CfbCipher_Encrypt(t *testing.T) {
	cipher := NewAes256CfbCipher("1234567890123456")
	bf := []byte("123435bhfdbsjhgbhadbcbxjhcbxzjhc")
	encrypt, _ := cipher.Encrypt(bf)
	fmt.Printf("encode %s \n", base64.StdEncoding.EncodeToString(encrypt))
	//decodeString, _ := base64.StdEncoding.DecodeString("cwkj9E+Es8nUqd+7K6CXhSSRlqqXU+I+t2N2qFv0J6WMBjRad48p4oGlrZzyeU/5")
	bff, _ := cipher.Decrypt(encrypt)
	fmt.Printf("decode %s \n", string(bff))

	//decodeString, _ := base64.StdEncoding.DecodeString("5I+LHPPBeHyyShy5WeEwzb++mfguQQ==")
	//b, _ := cipher.Decrypt(decodeString)
	//fmt.Printf("%s \n", string(b))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//decrypt, _ := cipher.Decrypt(encrypt)
	//println(string(decrypt))
}

func TestParseGet(t *testing.T) {
	cipher := NewAes256CfbCipher("1234567890123456")
	decodeString, _ := base64.StdEncoding.DecodeString("5I+LHPPBeHyyShy5WeEwzb++mfguQQ==")
	b, _ := cipher.Decrypt(decodeString)
	fmt.Printf("%s \n", string(b))
}
