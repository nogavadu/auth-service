package user

import (
	"context"
	"fmt"
	"github.com/nogavadu/auth-service/internal/domain/model"
	"github.com/nogavadu/auth-service/internal/repository"
	roleRepoModel "github.com/nogavadu/auth-service/internal/repository/role/model"
	userRepoModel "github.com/nogavadu/auth-service/internal/repository/user/model"
	"github.com/nogavadu/auth-service/internal/service"
	"github.com/nogavadu/platform_common/pkg/db"
	"log/slog"
)

type userService struct {
	log *slog.Logger

	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	txManager db.TxManager
}

func New(
	log *slog.Logger,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	txManager db.TxManager,
) service.UserService {
	return &userService{
		log:       log,
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		txManager: txManager,
	}
}

func (s *userService) GetById(ctx context.Context, id int) (*model.User, error) {
	const op = "userService.GetById"
	log := s.log.With(slog.String("op", op))

	var user *userRepoModel.User
	var role string
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to get user by id", slog.String("error", errTx.Error()))
			}
		}()

		user, errTx = s.userRepo.GetById(ctx, id)
		if errTx != nil {
			return errTx
		}

		repoRole, errTx := s.roleRepo.GetById(ctx, user.RoleId)
		if errTx != nil {
			return errTx
		}

		role = repoRole.Name

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.User{
		Id: user.Id,
		UserInfo: model.UserInfo{
			Name:   user.Name,
			Email:  user.Email,
			Avatar: user.Avatar,
			Role:   role,
		},
	}, nil
}

func (s *userService) Update(ctx context.Context, id int, input *model.UserUpdateInput) error {
	const op = "userService.Update"
	log := s.log.With(slog.String("op", op))

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to get user by id", slog.String("error", errTx.Error()))
			}
		}()

		var role *roleRepoModel.Role
		if input.Role != nil {
			role, errTx = s.roleRepo.GetByName(ctx, *input.Role)
			if errTx != nil {
				return errTx
			}
		}
		roleId := int(role.ID)

		if errTx = s.userRepo.Update(ctx, id, &userRepoModel.UserUpdateInput{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
			Avatar:   input.Avatar,
			RoleId:   &roleId,
		}); errTx != nil {
			return errTx
		}

		return nil
	})

	return err
}

func (s *userService) Delete(ctx context.Context, id int) error {
	const op = "userService.Delete"

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
