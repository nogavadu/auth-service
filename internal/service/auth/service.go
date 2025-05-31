package auth

import (
	"context"
	"errors"
	"github.com/nogavadu/auth-service/internal/domain/model"
	"github.com/nogavadu/auth-service/internal/repository"
	userRepoModel "github.com/nogavadu/auth-service/internal/repository/user/model"
	"github.com/nogavadu/auth-service/internal/service"
	"github.com/nogavadu/auth-service/internal/utils"
	"github.com/nogavadu/platform_common/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrAlreadyExists       = errors.New("already exists")
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

	userRepo  repository.UserRepository
	txManager db.TxManager
}

func New(
	log *slog.Logger,
	refreshTokenSecret string,
	refreshTokenExp time.Duration,
	accessTokenSecret string,
	accessTokenExp time.Duration,
	userRepo repository.UserRepository,
	txManager db.TxManager,
) service.AuthService {
	return &authService{
		log:                 log,
		refreshTokenSecret:  refreshTokenSecret,
		refreshTokenExpTime: refreshTokenExp,
		accessTokenSecret:   accessTokenSecret,
		accessTokenExpTime:  accessTokenExp,
		userRepo:            userRepo,
		txManager:           txManager,
	}
}

func (s *authService) Register(ctx context.Context, userInfo *model.UserInfo, password string) (int, error) {
	const op = "authService.Register"

	log := s.log.With(slog.String("op", op))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))

		return 0, ErrInvalidCredentials
	}

	userId, err := s.userRepo.Create(ctx, &userRepoModel.UserInfo{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		PassHash: string(passHash),
		Avatar:   nil,
		RoleId:   0,
	}, string(passHash))
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return 0, ErrAlreadyExists
		}

		log.Error("failed to create user", slog.String("error", err.Error()))
		return 0, ErrInternal
	}

	return userId, nil
}

func (s *authService) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "authService.Login"

	log := s.log.With(slog.String("op", op))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}

		log.Error("failed to get user", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	if !utils.VerifyPassword(user.PassHash, password) {
		return "", ErrInvalidCredentials
	}

	refreshToken, err := utils.GenerateToken(
		&model.User{
			Id: user.Id,
			UserInfo: model.UserInfo{
				Email:  user.Email,
				RoleId: user.RoleId,
			},
		},
		s.refreshTokenSecret,
		s.refreshTokenExpTime,
	)
	if err != nil {
		log.Error("failed to generate jwt token", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	return refreshToken, nil
}

func (s *authService) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	const op = "authService.GetRefreshToken"
	log := s.log.With(slog.String("op", op))

	claims, err := utils.VerifyToken(refreshToken, s.refreshTokenSecret)
	if err != nil {
		return "", ErrInvalidRefreshToken
	}

	user, err := s.userRepo.GetByEmail(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}

		log.Error("failed to get user", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	newRefreshToken, err := utils.GenerateToken(
		&model.User{
			Id: user.Id,
			UserInfo: model.UserInfo{
				Email:  user.Email,
				RoleId: user.RoleId,
			},
		},
		s.refreshTokenSecret,
		s.refreshTokenExpTime,
	)
	if err != nil {
		log.Error("failed to generate jwt token", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	return newRefreshToken, nil
}

func (s *authService) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	const op = "authService.GetAccessToken"
	log := s.log.With(slog.String("op", op))

	claims, err := utils.VerifyToken(refreshToken, s.refreshTokenSecret)
	if err != nil {
		return "", ErrInvalidRefreshToken
	}

	user, err := s.userRepo.GetByEmail(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}

		log.Error("failed to get user", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	accessToken, err := utils.GenerateToken(
		&model.User{
			Id: user.Id,
			UserInfo: model.UserInfo{
				Email:  user.Email,
				RoleId: user.RoleId,
			},
		},
		s.accessTokenSecret,
		s.accessTokenExpTime,
	)
	if err != nil {
		log.Error("failed to generate jwt token", slog.String("err", err.Error()))
	}

	return accessToken, nil
}
