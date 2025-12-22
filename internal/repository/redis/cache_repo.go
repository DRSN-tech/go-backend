package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DRSN-tech/go-backend/internal/cfg"
	"github.com/DRSN-tech/go-backend/internal/repository/redis/converter"
	"github.com/DRSN-tech/go-backend/internal/usecase"
	"github.com/DRSN-tech/go-backend/pkg/clients"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/DRSN-tech/go-backend/pkg/logger"
	"github.com/jimlawless/whereami"
)

type CacheRepo struct {
	client *clients.RedisClient
	conv   converter.ProductInfoConverter
	cfg    *cfg.RedisCfg
	logger logger.Logger
}

func NewCacheRepo(client *clients.RedisClient, conv converter.ProductInfoConverter,
	cfg *cfg.RedisCfg, logger logger.Logger) *CacheRepo {
	return &CacheRepo{
		client: client,
		conv:   conv,
		cfg:    cfg,
		logger: logger,
	}
}

// GetProducts возвращает закэшированные продукты по ID, игнорируя промахи и логируя их
func (r *CacheRepo) GetProducts(ctx context.Context, ids []int64) (map[int64]usecase.ProductInfo, error) {
	keys := r.buildProductCacheKeys(ids)

	values, err := r.client.Client.MGet(ctx, keys...).Result()
	if err != nil {
		r.logger.Warnf("Redis MGET failed: %v", e.Wrap(whereami.WhereAmI(), err))
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	result := make(map[int64]usecase.ProductInfo, len(values))
	for i, val := range values {
		data, err := redisValueToBytes(val, keys[i])
		if err != nil {
			r.logger.Warnf("%v", e.Wrap(whereami.WhereAmI(), err))
		}

		if data == nil {
			continue // cache miss
		}

		model, err := r.unmarshalProductFromCache(data)
		if err != nil {
			r.logger.Warnf("Redis unmarshal failed: %v", e.Wrap(whereami.WhereAmI(), err))
			continue
		}

		if model.ID != ids[i] {
			r.logger.Warnf("Cache ID mismatch: key_id: %d, model_id: %d", ids[i], model.ID)
			if err := r.client.Client.Del(context.Background(), keys[i]).Err(); err != nil {
				r.logger.Warnf("Redis del failed: %v", e.Wrap(whereami.WhereAmI(), err))
			}
			continue // cache miss
		}
		result[ids[i]] = *r.conv.ToUseCase(model)
	}

	return result, nil
}

// SetProducts атомарно кэширует несколько продуктов с заданным TTL.
// Игнорирует ошибки сериализации/записи, логируя их.
func (r *CacheRepo) SetProducts(ctx context.Context, products []usecase.ProductInfo) error {
	models := r.conv.ToArrRedisModel(products)

	pipeline := r.client.Client.Pipeline()
	for _, model := range models {
		data, err := r.marshalProductForCache(model)
		if err != nil {
			r.logger.Warnf("Failed to marshal product for caching (Product ID: %d): %v", model.ID, e.Wrap(whereami.WhereAmI(), err))
			continue
		}

		key := r.productKey(model.ID)
		pipeline.Set(ctx, key, data, r.cfg.ProductTTL)
	}

	if _, err := pipeline.Exec(ctx); err != nil {
		r.logger.Warnf("Cache pipeline failed: %v", e.Wrap(whereami.WhereAmI(), err))
	}

	return nil
}

// DeleteProducts удаляет продукты из кэша по ID
func (r *CacheRepo) DeleteProducts(ctx context.Context, ids []int64) error {
	keys := r.buildProductCacheKeys(ids)

	if err := r.client.Client.Del(ctx, keys...).Err(); err != nil {
		r.logger.Warnf("Redis DEL failed: %v", e.Wrap(whereami.WhereAmI(), err))
	}

	return nil
}

// marshalProductForCache сериализует продукт в JSON для кэша
func (r *CacheRepo) marshalProductForCache(model converter.ProductInfoRedisModel) ([]byte, error) {
	data, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// unmarshalProductFromCache десериализует JSON из кэша в модель продукта
func (r *CacheRepo) unmarshalProductFromCache(data []byte) (*converter.ProductInfoRedisModel, error) {
	var model converter.ProductInfoRedisModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

// buildProductCacheKeys формирует Redis-ключи из ID продуктов
func (r *CacheRepo) buildProductCacheKeys(ids []int64) []string {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = r.productKey(id)
	}

	return keys
}

// productKey возвращает Redis-ключ для одного продукта
func (r *CacheRepo) productKey(id int64) string {
	return fmt.Sprintf("product:%d", id)
}

// redisValueToBytes конвертирует значение из Redis в []byte.
// Поддерживает string и []byte, возвращает ошибку для неизвестных типов.
func redisValueToBytes(val interface{}, key string) ([]byte, error) {
	switch v := val.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case nil:
		return nil, nil // cache miss
	default:
		return nil, fmt.Errorf("unexpected Redis value type for key %s: %T", key, val)
	}
}
