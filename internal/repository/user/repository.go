package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/nogavadu/auth-service/internal/repository"
	userRepoModel "github.com/nogavadu/auth-service/internal/repository/user/model"
)

type userRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repo.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, email string, passHash string) (uint64, error) {
	const op = "userRepository.Create"

	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var id uint64
	if err := r.db.QueryRow(ctx, query, email, passHash).Scan(&id); err != nil {
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

	query := `
		SELECT (id, email, password_hash, role)
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var user userRepoModel.User
	if err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PassHash, &user.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repo.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
