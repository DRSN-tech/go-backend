package pgdb

import (
	"context"

	"github.com/DRSN-tech/go-backend/internal/domain"
	"github.com/DRSN-tech/go-backend/internal/repository/pgdb/converter"
	"github.com/DRSN-tech/go-backend/internal/usecase"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/DRSN-tech/go-backend/pkg/tr"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jimlawless/whereami"
)

// ProductRepo реализует репозиторий продуктов поверх PostgreSQL.
type ProductRepo struct {
	pool *pgxpool.Pool
	conv converter.ProductConverter
}

func NewProductRepo(pool *pgxpool.Pool, conv converter.ProductConverter) *ProductRepo {
	return &ProductRepo{
		pool: pool,
		conv: conv,
	}
}

// Upsert идемпотентно создаёт или обновляет продукт по уникальному имени,
// Запись обновляется только при изменении цены или категории.
func (p *ProductRepo) Upsert(ctx context.Context, product *domain.Product) (*usecase.UpsertProductRes, error) {
	tx, err := tr.TxFromCtx(ctx)
	if err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	// VALUES ($1, $2, $3) name, price, category_id
	query := `
		WITH upsert AS (
		INSERT INTO products (name, price, category_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (name)
		DO UPDATE SET
			price = EXCLUDED.price,
			category_id = EXCLUDED.category_id,
			updated_at = NOW()
		WHERE
			products.price IS DISTINCT FROM EXCLUDED.price OR
			products.category_id IS DISTINCT FROM EXCLUDED.category_id
		RETURNING
			id, name, price, category_id, created_at, updated_at, is_archived
		)
		SELECT
			id, name, price, category_id, created_at, updated_at, is_archived,
			false AS no_changes
		FROM upsert
		
		UNION ALL
		
		SELECT
			id, name, price, category_id, created_at, updated_at, is_archived,
			true AS no_changes
		FROM products
		WHERE name = $1
		  AND NOT EXISTS (SELECT 1 FROM upsert);
	`

	var model converter.ProductModel
	var noChanges bool
	err = tx.QueryRow(ctx, query, product.Name, product.Price, product.CategoryID).
		Scan(
			&model.ID, &model.Name, &model.Price, &model.CategoryID,
			&model.CreatedAt, &model.UpdatedAt, &model.IsArchived, &noChanges,
		)
	if err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	return usecase.NewUpsertProductRes(p.conv.ToEntity(&model), noChanges), nil
}

// GetProductsInfo возвращает информацию о продуктах по их идентификаторам, включая название категории.
func (p *ProductRepo) GetProductsInfo(ctx context.Context, ids []int64) ([]usecase.ProductInfo, error) {
	query := `
		SELECT pr.id, pr.name, pr.price, cat.name
		FROM products pr
		JOIN categories cat ON pr.category_id = cat.id
		WHERE pr.id = ANY($1)
	`

	rows, err := p.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}
	defer rows.Close()

	result := make([]usecase.ProductInfo, 0)
	for rows.Next() {
		var product usecase.ProductInfo
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.CategoryName); err != nil {
			return nil, e.Wrap(whereami.WhereAmI(), err)
		}

		result = append(result, product)
	}

	return result, nil
}
