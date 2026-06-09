package jwtx

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/loanem-backend/auth-service/config"
	"github.com/loanem-backend/auth-service/internal/entity"
)

func GenerateAccessToken(a *entity.Assistant) (string, error) {
	timeNow := time.Now()

	expireDurationStr := config.GetEnv("JWT_ACCESS_EXPIRE", "30m")
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

	privateKey, err := loadPrivate(config.GetEnv("JWT_PRIVATE_KEY", "./keys/jwt_private.pem"))
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}

func GenerateRefreshToken(aID int) (time.Duration, string, error) {
	timeNow := time.Now()

	expireDurationStr := config.GetEnv("JWT_REFRESH_EXPIRE", "36h")
	expDuration, err := time.ParseDuration(expireDurationStr)
	if err != nil {
		return 0, "", err
	}

	claims := Claims{
		ID: aID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth-service",
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(expDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := loadPrivate(config.GetEnv("JWT_PRIVATE_KEY", "./keys/jwt_private.pem"))
	if err != nil {
		return 0, "", err
	}

	signedToken, err := token.SignedString(privateKey)

	return expDuration, signedToken, err
}
