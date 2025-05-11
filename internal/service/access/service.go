package access

import (
	"context"
	"errors"
	"github.com/nogavadu/auth-service/internal/repository"
	"github.com/nogavadu/auth-service/internal/service"
	"github.com/nogavadu/auth-service/internal/utils"
	"log/slog"
	"time"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrPermissionDenied = errors.New("access denied")
	ErrInternal         = errors.New("internal error")
)

type accessService struct {
	log *slog.Logger

	refreshTokenSecret  string
	refreshTokenExpTime time.Duration
	accessTokenSecret   string
	accessTokenExpTime  time.Duration

	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func New(
	log *slog.Logger,
	refreshTokenSecret string,
	refreshTokenExp time.Duration,
	accessTokenSecret string,
	accessTokenExp time.Duration,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) service.AccessService {
	return &accessService{
		log:                 log,
		refreshTokenSecret:  refreshTokenSecret,
		refreshTokenExpTime: refreshTokenExp,
		accessTokenSecret:   accessTokenSecret,
		accessTokenExpTime:  accessTokenExp,
		userRepo:            userRepo,
		roleRepo:            roleRepo,
	}
}

func (s *accessService) Check(ctx context.Context, accessToken string, requiredLvl int) error {
	const op = "accessService.Check"

	claims, err := utils.VerifyToken(accessToken, s.accessTokenSecret)
	if err != nil {
		return ErrInvalidToken
	}

	role, err := s.roleRepo.GetByName(ctx, claims.Role)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrInvalidToken
		}

		s.log.Error("%s: %w", op, err)
		return ErrInternal
	}

	if role.Level < requiredLvl {
		return ErrPermissionDenied
	}

	return nil
}
