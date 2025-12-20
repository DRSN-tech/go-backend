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

func EnsureCollection(ctx context.Context, client *QdrantClient) error {
	exists, err := client.Client.CollectionExists(ctx, client.cfg.QdrantCollectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	if !exists {
		if err := client.Client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: client.cfg.QdrantCollectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     client.cfg.VectorSize,
				Distance: qdrant.Distance_Cosine,
			}),
		}); err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
	}

	return nil
}
