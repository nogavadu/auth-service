package service

import (
	"context"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (int, error)
	Login(ctx context.Context, email string, password string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, accessToken string, requiredLvl int) error
}
