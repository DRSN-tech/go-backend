package clients

import (
	"context"
	"fmt"

	config "github.com/DRSN-tech/go-backend/internal/cfg"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/jimlawless/whereami"
	"github.com/qdrant/go-client/qdrant"
)

type QdrantClient struct {
	Client *qdrant.Client
	cfg    *config.QdrantCfg
}

func NewQdrantClient(cfg *config.QdrantCfg) (*QdrantClient, error) {
	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host:   cfg.Host,
		Port:   cfg.Port,
		APIKey: cfg.ApiKey,
		UseTLS: cfg.UseTLS,
	})
	if err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	return &QdrantClient{
		Client: qdrantClient,
		cfg:    cfg,
	}, nil
}

func (q *QdrantClient) EnsureCollection(ctx context.Context) error {
	exists, err := q.Client.CollectionExists(ctx, q.cfg.QdrantCollectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	if !exists {
		if err := q.Client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: q.cfg.QdrantCollectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     q.cfg.VectorSize,
				Distance: qdrant.Distance_Cosine,
			}),
		}); err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
	}

	return nil
}
