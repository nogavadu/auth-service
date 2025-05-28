package user

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	repo "github.com/nogavadu/auth-service/internal/repository"
	userRepoModel "github.com/nogavadu/auth-service/internal/repository/user/model"
	"github.com/nogavadu/platform_common/pkg/db"
)

type userRepository struct {
	dbc db.Client
}

func New(dbc db.Client) repo.UserRepository {
	return &userRepository{
		dbc: dbc,
	}
}

func (r *userRepository) Create(ctx context.Context, email string, passHash string) (int, error) {
	const op = "userRepository.Create"

	queryRaw, args, err := sq.
		Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("email", "password_hash").
		Values(email, passHash).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query := db.Query{
		Name:     "userRepository.Create",
		QueryRaw: queryRaw,
	}

	var id int
	if err = r.dbc.DB().ScanOneContext(ctx, &id, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == repo.PgErrAlreadyExistsCode {
			return 0, fmt.Errorf("%s: %w", op, repo.ErrAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*userRepoModel.User, error) {
	const op = "userRepository.GetByEmail"

	queryRaw, args, err := sq.
		Select("id", "email", "password_hash", "role").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query := db.Query{
		Name:     "userRepository.GetByEmail",
		QueryRaw: queryRaw,
	}

	var user userRepoModel.User
	if err = r.dbc.DB().ScanOneContext(ctx, &user, query, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repo.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
