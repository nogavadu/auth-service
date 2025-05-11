package role

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/nogavadu/auth-service/internal/repository"
	roleRepoModel "github.com/nogavadu/auth-service/internal/repository/role/model"
)

type roleRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repo.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*roleRepoModel.Role, error) {
	const op = "roleRepository.GetByName"

	query := `
		SELECT (id, name, level)
		FROM roles
		WHERE name = $1
		LIMIT 1
	`

	var role roleRepoModel.Role
	if err := r.db.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name, &role.Level); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repo.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}
