package minio

import (
	"bytes"
	"context"

	"github.com/DRSN-tech/go-backend/internal/cfg"
	"github.com/DRSN-tech/go-backend/internal/domain"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/jimlawless/whereami"
	"github.com/minio/minio-go/v7"
)

// ImageRepo реализует репозиторий изображений поверх MinIO.
type ImageRepo struct {
	mc  *minio.Client
	cfg *cfg.MinIOCfg
}

func NewImageRepo(mc *minio.Client, cfg *cfg.MinIOCfg) *ImageRepo {
	return &ImageRepo{
		mc:  mc,
		cfg: cfg,
	}
}

// Upload загружает изображение в MinIO и возвращает ключ объекта.
func (i *ImageRepo) Upload(ctx context.Context, image *domain.Image) (string, error) {
	reader := bytes.NewReader(image.Bytes)

	info, err := i.mc.PutObject(ctx, i.cfg.BucketName, image.ObjectKey, reader, *image.Size, minio.PutObjectOptions{
		ContentType: *image.MimeType,
	})
	if err != nil {
		return "", e.Wrap(whereami.WhereAmI(), err)
	}

	return info.Key, nil
}

// Delete удаляет объект из MinIO по указанному ключу.
func (i *ImageRepo) Delete(ctx context.Context, key string) error {
	if err := i.mc.RemoveObject(ctx, i.cfg.BucketName, key, minio.RemoveObjectOptions{}); err != nil {
		return e.Wrap(whereami.WhereAmI(), err)
	}

	return nil
}
