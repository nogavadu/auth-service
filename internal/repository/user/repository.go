package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
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

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*userRepoModel.User, error) {
	const op = "userRepo.GetByEmail"

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
