package domain

import "time"

type ProductEmbeddingVersion struct {
	ID               int64
	ProductID        int64
	EmbeddingVersion int32
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	IsArchived       bool
}
