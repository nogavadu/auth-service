package auth

import (
	"context"
	"errors"
	"github.com/nogavadu/auth-service/internal/domain/model"
	"github.com/nogavadu/auth-service/internal/repository"
	"github.com/nogavadu/auth-service/internal/service"
	"github.com/nogavadu/auth-service/internal/utils"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInternal            = errors.New("internal error")
)

type authService struct {
	log *slog.Logger

	refreshTokenSecret  string
	refreshTokenExpTime time.Duration
	accessTokenSecret   string
	accessTokenExpTime  time.Duration

	repo repository.UserRepository
}

func New(
	log *slog.Logger,
	refreshTokenSecret string,
	refreshTokenExp time.Duration,
	accessTokenSecret string,
	accessTokenExp time.Duration,
	repo repository.UserRepository,
) service.AuthService {
	return &authService{
		log:                 log,
		refreshTokenSecret:  refreshTokenSecret,
		refreshTokenExpTime: refreshTokenExp,
		accessTokenSecret:   accessTokenSecret,
		accessTokenExpTime:  accessTokenExp,
		repo:                repo,
	}
}

func (s *authService) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "authService.Login"

	s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}

		s.log.Error("failed to get user", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	if !utils.VerifyPassword(user.PassHash, password) {
		return "", ErrInvalidCredentials
	}

	refreshToken, err := utils.GenerateToken(
		&model.UserInfo{
			Email: user.Email,
			Role:  user.Role,
		},
		s.refreshTokenSecret,
		s.refreshTokenExpTime,
	)
	if err != nil {
		s.log.Error("failed to generate jwt token", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	return refreshToken, nil
}

func (s *authService) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	const op = "authService.GetRefreshToken"
	s.log.With(slog.String("op", op))

	claims, err := utils.VerifyToken(refreshToken, s.refreshTokenSecret)
	if err != nil {
		return "", ErrInvalidRefreshToken
	}

	user, err := s.repo.GetByEmail(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}

		s.log.Error("failed to get user", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	newRefreshToken, err := utils.GenerateToken(
		&model.UserInfo{
			Email: user.Email,
			Role:  user.Role,
		},
		s.refreshTokenSecret,
		s.refreshTokenExpTime,
	)
	if err != nil {
		s.log.Error("failed to generate jwt token", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	return newRefreshToken, nil
}

func (s *authService) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	const op = "authService.GetAccessToken"
	s.log.With(slog.String("op", op))

	claims, err := utils.VerifyToken(refreshToken, s.refreshTokenSecret)
	if err != nil {
		return "", ErrInvalidRefreshToken
	}

	user, err := s.repo.GetByEmail(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}

		s.log.Error("failed to get user", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	accessToken, err := utils.GenerateToken(
		&model.UserInfo{
			Email: user.Email,
			Role:  user.Role,
		},
		s.accessTokenSecret,
		s.accessTokenExpTime,
	)
	if err != nil {
		s.log.Error("failed to generate jwt token", slog.String("err", err.Error()))
	}

	return accessToken, nil
}
