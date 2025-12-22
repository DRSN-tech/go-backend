package clients

import (
	"context"

	"github.com/DRSN-tech/go-backend/internal/cfg"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/jimlawless/whereami"
	r "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *r.Client
}

func NewRedisClient(cfg *cfg.RedisCfg) *RedisClient {
	client := r.NewClient(&r.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	return &RedisClient{
		Client: client,
	}
}

func (r *RedisClient) Ping(ctx context.Context) error {
	if err := r.Client.Ping(ctx).Err(); err != nil {
		// r.logger.Errorf(err, "failed to connect to redis server: %s\n", err.Error())
		return e.Wrap(whereami.WhereAmI(), err)
	}

	return nil
}
