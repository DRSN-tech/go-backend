package pgdb

import (
	"context"

	"github.com/DRSN-tech/go-backend/internal/domain"
	"github.com/DRSN-tech/go-backend/internal/repository/pgdb/converter"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/DRSN-tech/go-backend/pkg/tr"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jimlawless/whereami"
)

// CategoryRepo реализует репозиторий категорий поверх PostgreSQL.
type CategoryRepo struct {
	pool *pgxpool.Pool
	conv converter.CategoryConverter
}

func NewCategoryRepo(pool *pgxpool.Pool, conv converter.CategoryConverter) *CategoryRepo {
	return &CategoryRepo{pool: pool, conv: conv}
}

// Create идемпотентно создаёт категорию по имени, игнорируя дубликаты.
func (c *CategoryRepo) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	tx, err := tr.TxFromCtx(ctx)
	if err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	query := `
		INSERT INTO categories(name) VALUES ($1)
		ON CONFLICT (name) DO NOTHING
		RETURNING id, name, created_at, updated_at, is_archived;
	`

	var model converter.CategoryModel
	if err := tx.QueryRow(ctx, query, category.Name).
		Scan(
			&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt, &model.IsArchived,
		); err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	return c.conv.ToEntity(&model), nil
}
