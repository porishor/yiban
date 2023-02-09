// 从网页中获取 key ,并且加密密码
package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func EncryptPassword(public_key string, data []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(public_key))
	if block == nil {
		return nil, errors.New("public key error")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("parse public key error：%s", err.Error()))
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), data)
}
