package model

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims
	Id    int    `json:"id"`
	Email string `json:"Email"`
	Role  string `json:"role"`
}
