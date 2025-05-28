package role

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	repo "github.com/nogavadu/auth-service/internal/repository"
	roleRepoModel "github.com/nogavadu/auth-service/internal/repository/role/model"
	"github.com/nogavadu/platform_common/pkg/db"
)

type roleRepository struct {
	dbc db.Client
}

func New(dbc db.Client) repo.RoleRepository {
	return &roleRepository{
		dbc: dbc,
	}
}

func (r *roleRepository) GetById(ctx context.Context, id int) (*roleRepoModel.Role, error) {
	const op = "roleRepository.GetByName"

	queryRaw, args, err := sq.
		Select("id", "name", "level").
		PlaceholderFormat(sq.Dollar).
		From("roles").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query := db.Query{
		Name:     "roleRepository.GetByName",
		QueryRaw: queryRaw,
	}

	var role roleRepoModel.Role
	if err = r.dbc.DB().ScanOneContext(ctx, &role, query, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repo.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}
