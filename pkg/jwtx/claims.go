package jwtx

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	JwtID string `json:"jwt_id"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
}
