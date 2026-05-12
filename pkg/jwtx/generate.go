package jwtx

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/loanem-backend/auth-service/config"
	"github.com/loanem-backend/auth-service/internal/entity"
)

func GenerateToken(a *entity.Assistant) (string, error) {
	timeNow := time.Now()

	expireDurationStr := config.GetEnv("JWT_EXPIRE", "30m")
	expDuration, err := time.ParseDuration(expireDurationStr)
	if err != nil {
		return "", err
	}

	claims := Claims{
		ID:   a.ID,
		Name: a.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth-service",
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(expDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := loadPrivate(config.GetEnv("JWT_PRIVATE_KEY", "./backend-infra/keys/jwt_private.pem"))
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}
