package qdrant

import (
	"context"

	"github.com/DRSN-tech/go-backend/internal/cfg"
	"github.com/DRSN-tech/go-backend/internal/domain"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/jimlawless/whereami"
	"github.com/qdrant/go-client/qdrant"
)

// EmbeddingRepo репозиторий для работы с embedding-векторами в Qdrant
type EmbeddingRepo struct {
	client *qdrant.Client
	cfg    *cfg.QdrantCfg
}

func NewEmbeddingRepo(client *qdrant.Client, cfg *cfg.QdrantCfg) *EmbeddingRepo {
	return &EmbeddingRepo{
		client: client,
		cfg:    cfg,
	}
}

// Upsert сохраняет или обновляет embedding-векторы в указанной коллекции Qdrant.
func (q *EmbeddingRepo) Upsert(ctx context.Context, vectors []domain.Embedding) error {
	reqVectors := make([]*qdrant.PointStruct, 0, len(vectors))
	for _, vector := range vectors {
		reqVectors = append(reqVectors, &qdrant.PointStruct{
			Id:      qdrant.NewIDUUID(vector.ID),
			Vectors: qdrant.NewVectors(vector.Vector...),
			Payload: qdrant.NewValueMap(vector.Payload),
		})
	}

	_, err := q.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: q.cfg.QdrantCollectionName,
		Points:         reqVectors,
	})
	if err != nil {
		return e.Wrap(whereami.WhereAmI(), err)
	}

	return nil
}
