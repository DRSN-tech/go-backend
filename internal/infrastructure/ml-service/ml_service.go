package ml_service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DRSN-tech/go-backend/internal/proto"
	"github.com/DRSN-tech/go-backend/internal/usecase"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/DRSN-tech/go-backend/pkg/jitter"
	"github.com/DRSN-tech/go-backend/pkg/logger"
)

// MLService клиент для взаимодействия с внешним ML-сервисом
type MLService struct {
	client        proto.MachineLearningServiceClient
	maxConcurrent int // TODO: 5-10
	maxRetries    int // TODO: 3
	logger        logger.Logger
}

func NewMLService(client proto.MachineLearningServiceClient, maxConcurrent int, maxRetries int, logger logger.Logger) *MLService {
	return &MLService{
		client:        client,
		maxConcurrent: maxConcurrent,
		maxRetries:    maxRetries,
		logger:        logger,
	}
}

// VectorizeRequest выполняет векторизацию изображений с retry-логикой и экспоненциальной задержкой
func (m *MLService) VectorizeRequest(ctx context.Context, req *usecase.VectorizeReq) ([]usecase.VectorizeRes, error) {
	const (
		op         = "MLService.VectorizeRequest"
		baseJitter = 1 * time.Second
		maxJitter  = 30 * time.Second
	)

	for attempt := 0; attempt < m.maxRetries; attempt++ {
		vectors, err := m.vectorizeBatch(ctx, req)
		if err == nil {
			return vectors, nil
		}

		if attempt == m.maxRetries-1 {
			return nil, e.Wrap(op, fmt.Errorf("all %d attempts failed", m.maxRetries))
		}

		sleepTime := jitter.ExponentialBackoff(
			baseJitter,
			maxJitter,
			attempt,
			jitter.DefaultJitter,
		)

		m.logger.Warnf("vectorization failed, retrying in %v (attempt %d)", sleepTime, attempt+1)
		select {
		case <-time.After(sleepTime):
		case <-ctx.Done():
			return nil, e.Wrap(op, ctx.Err())
		}
	}

	return nil, e.Wrap(op, fmt.Errorf("unreachable"))
}

// vectorizeBatch отправляет батч изображений на векторизацию параллельно с ограничением конкурентности
func (m *MLService) vectorizeBatch(ctx context.Context, req *usecase.VectorizeReq) ([]usecase.VectorizeRes, error) {
	const op = "MLService.vectorizeBatch"

	vectorCh := make(chan usecase.VectorizeRes, len(req.Images))
	errCh := make(chan error, len(req.Images))
	sem := make(chan struct{}, m.maxConcurrent)

	var wg sync.WaitGroup
	for _, image := range req.Images {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			protoReq := proto.VectorizeRequest{
				ImageData: image.Data,
				ImageType: req.ImageType,
			}

			res, err := m.client.VectorizeImage(ctx, &protoReq)
			if err != nil {
				errCh <- err
				return
			}

			vectorCh <- *usecase.NewVectorizeRes(res.Vector, res.ModelVersion)
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(vectorCh)
	}()

	vectors := make([]usecase.VectorizeRes, 0, len(req.Images))
	for completed := 0; completed < len(req.Images); {
		select {
		case vector, ok := <-vectorCh:
			if ok {
				vectors = append(vectors, vector)
				completed++
			}
		case err, ok := <-errCh:
			if ok {
				return nil, e.Wrap(op, err)
			}
		case <-ctx.Done():
			return nil, e.Wrap(op, ctx.Err())
		}
	}

	return vectors, nil
}
