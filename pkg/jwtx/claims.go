package jwtx

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	ID   int    `json:"id"`
	Name string `json:"name"`
}
