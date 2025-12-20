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

type ProductEmbeddingVersionRepo struct {
	pool *pgxpool.Pool
	conv converter.ProductEmbeddingVersionConverter
}

func NewProductEmbeddingVersionRepo(pool *pgxpool.Pool, conv converter.ProductEmbeddingVersionConverter) *ProductEmbeddingVersionRepo {
	return &ProductEmbeddingVersionRepo{
		pool: pool,
		conv: conv,
	}
}

func (p *ProductEmbeddingVersionRepo) Upsert(ctx context.Context, productID int64) (*domain.ProductEmbeddingVersion, error) {
	tx, err := tr.TxFromCtx(ctx)
	if err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	var model converter.ProductEmbeddingVersionModel
	query := `
	INSERT INTO product_embedding_version (product_id)
    VALUES ($1)
    ON CONFLICT (product_id)
    DO UPDATE SET embedding_version = product_embedding_version.embedding_version + 1,
                  updated_at = NOW()
    RETURNING id, product_id, embedding_version, created_at, updated_at, is_archived;
	`

	err = tx.QueryRow(ctx, query, productID).Scan(
		&model.ID,
		&model.ProductID,
		&model.EmbeddingVersion,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.IsArchived,
	)
	if err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	return p.conv.ToEntity(&model), nil
}
