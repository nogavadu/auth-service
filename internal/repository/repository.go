package repository

import (
	"context"
	"errors"
	roleRepoModel "github.com/nogavadu/auth-service/internal/repository/role/model"
	userRepoModel "github.com/nogavadu/auth-service/internal/repository/user/model"
)

const (
	PgNotFoundCode = "23505"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrInternal      = errors.New("internal error")
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*userRepoModel.User, error)
}

type RoleRepository interface {
	GetByName(ctx context.Context, name string) (*roleRepoModel.Role, error)
}
