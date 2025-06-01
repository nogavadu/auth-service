package service

import (
	"context"
	"github.com/nogavadu/auth-service/internal/domain/model"
)

type AuthService interface {
	Register(ctx context.Context, userInfo *model.UserInfo, password string) (int, error)
	Login(ctx context.Context, email string, password string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	IsUser(ctx context.Context, userId int, refreshToken string) error
}

type AccessService interface {
	Check(ctx context.Context, accessToken string, requiredLvl int) error
}

type UserService interface {
	GetById(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, id int, input *model.UserUpdateInput) error
	Delete(ctx context.Context, id int) error
}
