package jwtx

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func loadPrivate(path string) (*rsa.PrivateKey, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	b, _ := pem.Decode(content)

	return x509.ParsePKCS1PrivateKey(b.Bytes)
}
