package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nogavadu/auth-service/internal/domain/model"
	"time"
)

func GenerateToken(info *model.UserInfo, secretKey string, dur time.Duration) (string, error) {
	claims := &model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(dur).Unix(),
		},
		Id:     info.Id,
		Email:  info.Email,
		RoleId: info.RoleId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenStr, secretKey string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("enexpected token signing method")
			}

			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
