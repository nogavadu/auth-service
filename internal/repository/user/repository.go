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

func (r *userRepository) Create(ctx context.Context, info *userRepoModel.UserInfo, passHash string) (int, error) {
	const op = "userRepository.Create"

	values := map[string]interface{}{
		"email":         info.Email,
		"password_hash": passHash,
	}

	if info.Name != nil {
		values["name"] = *info.Name
	}

	queryRaw, args, err := sq.
		Insert("users").
		PlaceholderFormat(sq.Dollar).
		SetMap(values).
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
		Select("id", "name", "email", "avatar", "password_hash", "role").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query := db.Query{
		Name:     op,
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

func (r *userRepository) GetById(ctx context.Context, id int) (*userRepoModel.User, error) {
	const op = "userRepository.GetById"

	queryRaw, args, err := sq.
		Select("id", "name", "email", "avatar", "password_hash", "role").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query := db.Query{
		Name:     op,
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

func (r *userRepository) Update(ctx context.Context, id int, input *userRepoModel.UserUpdateInput) error {
	const op = "userRepository.Update"

	values := map[string]interface{}{}
	if input.Email != nil {
		values["email"] = *input.Email
	}
	if input.Name != nil {
		values["name"] = *input.Name
	}
	if input.Avatar != nil {
		values["avatar"] = *input.Avatar
	}
	if input.Password != nil {
		values["password_hash"] = *input.Password
	}
	if input.RoleId != nil {
		values["role"] = *input.RoleId
	}

	fmt.Println(values)

	queryRaw, args, err := sq.
		Update("users").
		PlaceholderFormat(sq.Dollar).
		SetMap(values).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query := db.Query{
		Name:     op,
		QueryRaw: queryRaw,
	}

	_, err = r.dbc.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	const op = "userRepository.Delete"

	queryRaw, args, err := sq.
		Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query := db.Query{
		Name:     op,
		QueryRaw: queryRaw,
	}

	_, err = r.dbc.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
