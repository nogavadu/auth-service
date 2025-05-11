package model

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"Email"`
	Role  string `json:"role"`
}
