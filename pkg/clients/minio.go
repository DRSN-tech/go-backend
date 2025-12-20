package clients

import (
	"context"

	config "github.com/DRSN-tech/go-backend/internal/cfg"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/jimlawless/whereami"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOClient(cfg *config.Config) (*minio.Client, error) {
	minioCLient, err := minio.New(cfg.Minio.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.MinioRootUser, cfg.Minio.MinioRootPassword, ""),
		Secure: cfg.Minio.MinioUseSSL,
	})
	if err != nil {
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	return minioCLient, nil
}

func EnsureBucket(ctx context.Context, client *minio.Client, bucketName string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}

	return nil
}
