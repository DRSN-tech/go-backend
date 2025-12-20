package usecase

import (
	"context"

	"github.com/DRSN-tech/go-backend/internal/domain"
)

type ProductRepository interface {
	Upsert(ctx context.Context, product *domain.Product) (*UpsertProductRes, error)
	GetProductsInfo(ctx context.Context, ids []int64) ([]ProductInfo, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) (*domain.Category, error)
}

type ImageRepository interface {
	Upload(ctx context.Context, image *domain.Image) (string, error)
	Delete(ctx context.Context, key string) error
}

type EmbeddingRepository interface {
	Upsert(ctx context.Context, vectors []domain.Embedding) error
}

type CacheRepository interface {
	GetProducts(ctx context.Context, ids []int64) (map[int64]ProductInfo, error)
	SetProducts(ctx context.Context, products []ProductInfo) error
	DeleteProducts(ctx context.Context, ids []int64) error
}
