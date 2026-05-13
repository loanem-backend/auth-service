package jwtx

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/loanem-backend/auth-service/config"
)

func DecodeToken(signedToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.GetEnv("JWT_SECRET_KEY", "")), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed parsing token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token.Claims.(*Claims), nil
}
