package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
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
	roleRepo  repository.RoleRepository
	txManager db.TxManager

	registrationsProducer sarama.SyncProducer
}

func New(
	log *slog.Logger,
	refreshTokenSecret string,
	refreshTokenExp time.Duration,
	accessTokenSecret string,
	accessTokenExp time.Duration,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	txManager db.TxManager,
) service.AuthService {
	var addresses = []string{"kafka1:29091", "kafka2:29092"}

	producer, err := sarama.NewSyncProducer(addresses, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	return &authService{
		log:                   log,
		refreshTokenSecret:    refreshTokenSecret,
		refreshTokenExpTime:   refreshTokenExp,
		accessTokenSecret:     accessTokenSecret,
		accessTokenExpTime:    accessTokenExp,
		userRepo:              userRepo,
		roleRepo:              roleRepo,
		txManager:             txManager,
		registrationsProducer: producer,
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

	var userId int
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error(fmt.Sprintf("%s: failed transaction, %w", op, errTx))
			}
		}()

		rUserId, errTx := s.userRepo.Create(ctx, &userRepoModel.UserInfo{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			PassHash: string(passHash),
			Avatar:   nil,
			RoleId:   0,
		})
		if errTx != nil {
			if errors.Is(errTx, repository.ErrAlreadyExists) {
				return ErrAlreadyExists
			}

			return ErrInternal
		}
		userId = rUserId

		_, _, errTx = s.registrationsProducer.SendMessage(&sarama.ProducerMessage{
			Topic:     "registrations-topic",
			Value:     sarama.StringEncoder(userInfo.Email),
			Timestamp: time.Now(),
		})
		if errTx != nil {
			return ErrInternal
		}

		return nil
	})

	return userId, err
}

func (s *authService) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "authService.Login"

	log := s.log.With(slog.String("op", op))

	var user model.User
	var role string
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to get user", slog.String("error", errTx.Error()))
			}
		}()

		repoUser, errTx := s.userRepo.GetByEmail(ctx, email)
		if errTx != nil {
			if errors.Is(errTx, repository.ErrNotFound) {
				return ErrInvalidCredentials
			}

			return ErrInternal
		}

		repoRole, errTx := s.roleRepo.GetById(ctx, repoUser.RoleId)
		if errTx != nil {
			if errors.Is(errTx, repository.ErrNotFound) {
				return ErrInvalidCredentials
			}

			return ErrInternal
		}
		role = repoRole.Name

		user = model.User{
			Id: repoUser.Id,
			UserInfo: model.UserInfo{
				Name:   repoUser.Name,
				Email:  repoUser.Email,
				Avatar: repoUser.Avatar,
				Role:   role,
			},
		}

		if !utils.VerifyPassword(repoUser.PassHash, password) {
			return ErrInvalidCredentials
		}

		return nil
	})
	if err != nil {
		return "", ErrInternal
	}

	refreshToken, err := utils.GenerateToken(
		&model.User{
			Id: user.Id,
			UserInfo: model.UserInfo{
				Email: user.Email,
				Role:  role,
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

	var user model.User
	var role string
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to get user", slog.String("error", errTx.Error()))
			}
		}()

		repoUser, errTx := s.userRepo.GetById(ctx, claims.Id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrInvalidCredentials
			}

			return ErrInternal
		}

		repoRole, errTx := s.roleRepo.GetById(ctx, repoUser.RoleId)
		if errTx != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrInvalidCredentials
			}

			return ErrInternal
		}
		role = repoRole.Name

		user = model.User{
			Id: repoUser.Id,
			UserInfo: model.UserInfo{
				Name:   repoUser.Name,
				Email:  repoUser.Email,
				Avatar: repoUser.Avatar,
				Role:   role,
			},
		}

		return nil
	})
	if err != nil {
		return "", ErrInternal
	}

	newRefreshToken, err := utils.GenerateToken(
		&model.User{
			Id: user.Id,
			UserInfo: model.UserInfo{
				Email: user.Email,
				Role:  role,
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

	var user model.User
	var role string
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to get user", slog.String("error", errTx.Error()))
			}
		}()

		repoUser, errTx := s.userRepo.GetById(ctx, claims.Id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrInvalidCredentials
			}

			return ErrInternal
		}

		repoRole, errTx := s.roleRepo.GetById(ctx, repoUser.RoleId)
		if errTx != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrInvalidCredentials
			}

			return ErrInternal
		}
		role = repoRole.Name

		user = model.User{
			Id: repoUser.Id,
			UserInfo: model.UserInfo{
				Name:   repoUser.Name,
				Email:  repoUser.Email,
				Avatar: repoUser.Avatar,
				Role:   role,
			},
		}

		return nil
	})
	if err != nil {
		return "", ErrInternal
	}

	accessToken, err := utils.GenerateToken(
		&model.User{
			Id: user.Id,
			UserInfo: model.UserInfo{
				Email: user.Email,
				Role:  role,
			},
		},
		s.accessTokenSecret,
		s.accessTokenExpTime,
	)
	if err != nil {
		log.Error("failed to generate jwt token", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	return accessToken, nil
}

func (s *authService) IsUser(ctx context.Context, userId int, refreshToken string) error {
	const op = "authService.IsUser"
	log := s.log.With(slog.String("op", op))

	claims, err := utils.VerifyToken(refreshToken, s.refreshTokenSecret)
	if err != nil {
		log.Error("failed to verify jwt token", slog.String("err", err.Error()))
		return ErrInvalidRefreshToken
	}

	if claims.Id != userId {
		log.Error("not a same user")
		return ErrInvalidCredentials
	}

	return nil
}
