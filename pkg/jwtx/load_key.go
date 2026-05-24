package jwtx

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func loadPrivate(path string) (*rsa.PrivateKey, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(content)
}
