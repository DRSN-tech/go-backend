package converter

import (
	"time"

	"github.com/DRSN-tech/go-backend/internal/usecase"
	"github.com/google/uuid"
)

// ProductModel представляет запись таблицы product_types в PostgreSQL.
type ProductModel struct {
	ID         int64      `db:"id"`
	Name       string     `db:"name"`
	Price      int64      `db:"price"`
	CategoryID int64      `db:"category_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	IsArchived bool       `db:"is_archived"`
}

// CategoryModel представляет запись таблицы categories в PostgreSQL.
type CategoryModel struct {
	ID         int64      `db:"id"`
	Name       string     `db:"name"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	IsArchived bool       `db:"is_archived"`
}

// ProductEmbeddingVersionModel представляет запись таблицы product_embedding_version в PostgreSQL.
type ProductEmbeddingVersionModel struct {
	ID               int64      `db:"id"`
	ProductID        int64      `db:"product_id"`
	EmbeddingVersion int32      `db:"embedding_version"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
	IsArchived       bool       `db:"is_archived"`
}

type OutboxEventModel struct {
	ID                  int64                   `db:"id"`
	EventID             uuid.UUID               `db:"event_id"`
	ProductID           int64                   `db:"product_id"`
	EventType           usecase.OutboxEventType `db:"event_type"`
	Payload             []byte                  `db:"payload"`
	Status              usecase.OutboxStatus    `db:"status"`
	CreatedAt           time.Time               `db:"created_at"`
	ProcessingStartedAt *time.Time              `db:"processing_started_at"`
	ProcessedAt         *time.Time              `db:"processed_at"`
}
