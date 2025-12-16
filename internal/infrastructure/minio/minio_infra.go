package minio

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DRSN-tech/go-backend/internal/cfg"
	"github.com/DRSN-tech/go-backend/internal/domain"
	"github.com/DRSN-tech/go-backend/internal/infrastructure"
	"github.com/DRSN-tech/go-backend/internal/usecase"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/DRSN-tech/go-backend/pkg/logger"

	"github.com/google/uuid"
)

// MinioInfrastructure управляет загрузкой и очисткой изображений в MinIO.
type MinioInfrastructure struct {
	minioRepo         usecase.ImageRepository
	cfg               cfg.Configuration
	logger            logger.Logger
	shutdownCtx       context.Context
	wg                sync.WaitGroup
	uploadImagesLimit int
}

func NewMinioInfrastructure(minioRepo usecase.ImageRepository, cfg cfg.Configuration, logger logger.Logger, shutdownCtx context.Context) *MinioInfrastructure {
	return &MinioInfrastructure{
		minioRepo:         minioRepo,
		cfg:               cfg,
		logger:            logger,
		shutdownCtx:       shutdownCtx,
		wg:                sync.WaitGroup{},
		uploadImagesLimit: cfg.UploadImagesLimit,
	}
}

// UploadImages загружает изображения продукта в MinIO параллельно с ограничением одновременных операций.
// В случае ошибки отменяет остальные загрузки и запускает очистку уже загруженных файлов.
func (m *MinioInfrastructure) UploadImages(ctx context.Context, req *usecase.UploadImagesReq) (*usecase.UploadImagesRes, error) {
	const op = "MinioInfrastructure.UploadImages"
	// Отмена остальных загрузок при первой ошибке
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	keyCh := make(chan string, len(req.Images))
	errCh := make(chan error, len(req.Images))
	sem := make(chan struct{}, m.uploadImagesLimit)

	var uploadWg sync.WaitGroup
	for _, image := range req.Images {
		uploadWg.Add(1)
		go func() {
			defer uploadWg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			imageID := uuid.NewString()
			ext, err := infrastructure.GetExtensionFromMIME(image.MimeType)
			if err != nil {
				errCh <- fmt.Errorf("invalid mime type %s for %s: %w", image.MimeType, image.Name, err)
				return
			}
			objKey := fmt.Sprintf("%s/%s-%s.%s", req.Name, image.Name, imageID, ext)
			newImage := domain.NewImage(imageID, m.cfg.BucketName, objKey, image.Data, &image.Size, &image.MimeType)

			key, err := m.minioRepo.Upload(ctx, newImage)
			if err != nil {
				errCh <- fmt.Errorf("upload %s failed: %w", image.Name, err)
				return
			}

			keyCh <- key
		}()
	}

	go func() {
		uploadWg.Wait()
		close(errCh)
		close(keyCh)
	}()

	keys := make([]string, 0, len(req.Images))
	ok := false
	defer func() {
		if !ok && len(keys) > 0 {
			m.wg.Add(1)
			go m.cleanupUploadedKeys(keys)
		}
	}()

	for completed := 0; completed < len(req.Images); {
		select {
		case key, ok := <-keyCh:
			if ok {
				keys = append(keys, key)
				completed++
			}
		case err, ok := <-errCh:
			if ok {
				cancel()
				return nil, e.Wrap(op, err)
			}
		case <-ctx.Done():
			cancel()
			return nil, e.Wrap(op, ctx.Err())
		}
	}

	ok = true
	return usecase.NewUploadImagesRes(keys), nil
}

// CleanupImages запускает фоновую очистку указанных ключей MinIO
func (m *MinioInfrastructure) CleanupImages(keys []string) {
	if len(keys) == 0 {
		return
	}
	m.wg.Add(1)
	go m.cleanupUploadedKeys(keys)
}

// cleanupUploadedKeys удаляет указанные объекты из MinIO с экспоненциальной задержкой и jitter.
func (m *MinioInfrastructure) cleanupUploadedKeys(keys []string) {
	defer m.wg.Done() // сигнализируем завершение компенсации
	const op = "MinioInfrastructure.cleanupUploadedKeys"
	m.logger.Infof("%s: Cleaning up uploaded keys", op)

	// Создаём контекст с таймаутом на основе shutdownCtx
	ctx, cancel := context.WithTimeout(m.shutdownCtx, 30*time.Second)
	defer cancel()

	for _, key := range keys {
		backoff := time.Second
		for attempt := 0; attempt < 3; attempt++ {
			if err := m.minioRepo.Delete(ctx, key); err == nil {
				break // Успешно удалено
			}

			// Проверяем, не отменён ли контекст
			select {
			case <-ctx.Done():
				m.logger.Warnf("cleanup interrupted by shutdown, key=%v", key)
				return
			default:
			}

			if attempt < 2 {
				// Добавляем jitter для распределения нагрузки
				jitter := time.Duration(time.Now().UnixNano() % int64(time.Second))
				sleepTime := backoff + jitter

				select {
				case <-time.After(sleepTime):
				case <-ctx.Done():
					m.logger.Warnf("cleanup interrupted by shutdown during backoff, key=%v", key)
					return
				}
				backoff *= 2
			}
		}
	}
}

// WaitForCleanup ожидает завершения всех фоновых задач очистки с учётом таймаута завершения приложения.
func (m *MinioInfrastructure) WaitForCleanup(shutdownTimeoutCtx context.Context) error {
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-shutdownTimeoutCtx.Done():
		return fmt.Errorf("minio cleanup timeout during shutdown: %w", shutdownTimeoutCtx.Err())
	}
}
