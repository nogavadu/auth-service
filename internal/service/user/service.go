package user

import (
	"context"
	"fmt"
	"github.com/nogavadu/auth-service/internal/domain/model"
	"github.com/nogavadu/auth-service/internal/repository"
	userRepoModel "github.com/nogavadu/auth-service/internal/repository/user/model"
	"github.com/nogavadu/auth-service/internal/service"
	"github.com/nogavadu/platform_common/pkg/db"
	"log/slog"
)

type userService struct {
	log *slog.Logger

	userRepo  repository.UserRepository
	txManager db.TxManager
}

func New(log *slog.Logger, userRepo repository.UserRepository, txManager db.TxManager) service.UserService {
	return &userService{
		log:       log,
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (s *userService) GetById(ctx context.Context, id int) (*model.User, error) {
	const op = "userService.GetById"

	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.User{
		Id: user.Id,
		UserInfo: model.UserInfo{
			Name:   user.Name,
			Email:  user.Email,
			Avatar: user.Avatar,
			RoleId: user.RoleId,
		},
	}, nil
}

func (s *userService) Update(ctx context.Context, id int, input *model.UserUpdateInput) error {
	const op = "userService.Update"

	if err := s.userRepo.Update(ctx, id, &userRepoModel.UserUpdateInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Avatar:   input.Avatar,
		RoleId:   input.RoleId,
	}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	const op = "userService.Delete"

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
