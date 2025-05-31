package repository

import (
	"context"
	"errors"
	roleRepoModel "github.com/nogavadu/auth-service/internal/repository/role/model"
	userRepoModel "github.com/nogavadu/auth-service/internal/repository/user/model"
)

const (
	PgErrAlreadyExistsCode = "23505"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrInternal      = errors.New("internal error")
)

type UserRepository interface {
	Create(ctx context.Context, userInfo *userRepoModel.UserInfo, passHash string) (int, error)
	GetByEmail(ctx context.Context, email string) (*userRepoModel.User, error)
	GetById(ctx context.Context, id int) (*userRepoModel.User, error)
	Update(ctx context.Context, id int, input *userRepoModel.UserUpdateInput) error
	Delete(ctx context.Context, id int) error
}

type RoleRepository interface {
	GetByName(ctx context.Context, name string) (*roleRepoModel.Role, error)
	GetById(ctx context.Context, id int) (*roleRepoModel.Role, error)
}
