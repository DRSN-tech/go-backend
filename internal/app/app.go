package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/DRSN-tech/go-backend/internal/cfg"
	v1Grpc "github.com/DRSN-tech/go-backend/internal/delivery/v1/grpc"
	v1Http "github.com/DRSN-tech/go-backend/internal/delivery/v1/http"
	"github.com/DRSN-tech/go-backend/internal/infrastructure/kafka"
	minioInfra "github.com/DRSN-tech/go-backend/internal/infrastructure/minio"
	ml_service "github.com/DRSN-tech/go-backend/internal/infrastructure/ml-service"
	"github.com/DRSN-tech/go-backend/internal/proto"
	s3Repo "github.com/DRSN-tech/go-backend/internal/repository/minio"
	"github.com/DRSN-tech/go-backend/internal/repository/pgdb"
	pgdbConv "github.com/DRSN-tech/go-backend/internal/repository/pgdb/converter/generated"
	qdrantRepo "github.com/DRSN-tech/go-backend/internal/repository/qdrant"
	"github.com/DRSN-tech/go-backend/internal/repository/redis"
	redisConv "github.com/DRSN-tech/go-backend/internal/repository/redis/converter/generated"
	"github.com/DRSN-tech/go-backend/internal/usecase"
	"github.com/DRSN-tech/go-backend/pkg/clients"
	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/DRSN-tech/go-backend/pkg/logger"
	"github.com/DRSN-tech/go-backend/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/jimlawless/whereami"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FIXME: не пушить
func Run() {
	logger := logger.NewSlogLogger()

	cfg, err := config.Load(logger)
	if err != nil {
		logger.Errorf(err, "failed to load config")
		os.Exit(1)
	}

	logger.Infof("minio endpoint: %s", cfg.Minio.MinioEndpoint) // TODO: удалить

	db, err := initPGDB(logger, cfg)
	if err != nil {
		logger.Errorf(err, "failed to initialize database")
		os.Exit(1)
	}

	catConv := pgdbConv.NewCategoryConverterImpl()
	prConv := pgdbConv.NewProductConverterImpl()
	infoConv := redisConv.NewProductInfoConverterImpl()
	embConv := pgdbConv.NewProductEmbeddingVersionConverterImpl()

	productRepo := pgdb.NewProductRepo(db.Pool, prConv)
	categoryRepo := pgdb.NewCategoryRepo(db.Pool, catConv)
	prEmbeddingVersionRepo := pgdb.NewProductEmbeddingVersionRepo(db.Pool, embConv)

	minioClient, err := clients.NewMinIOClient(cfg)
	if err != nil {
		logger.Errorf(err, "failed to initialize minio client")
		os.Exit(1)
	}

	minioCtx, minioCancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := clients.EnsureBucket(minioCtx, minioClient, cfg.Minio.BucketName); err != nil {
		minioCancel()
		logger.Errorf(err, "failed to initialize MinIO bucket")
		os.Exit(1)
	}
	minioCancel()

	imageRepo := s3Repo.NewImageRepo(minioClient, cfg.Minio)

	qdrantClient, err := clients.NewQdrantClient(cfg.Qdrant)
	if err != nil {
		logger.Errorf(err, "failed to initialize qdrant")
		os.Exit(1)
	}
	qdrantCtx, qdrantCancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := clients.EnsureCollection(qdrantCtx, qdrantClient); err != nil {
		qdrantCancel()
		logger.Errorf(err, "failed to initialize qdrant")
		os.Exit(1)
	}
	qdrantCancel()

	embRepo := qdrantRepo.NewEmbeddingRepo(qdrantClient.Client, cfg.Qdrant)

	redisClient := clients.NewRedisClient(cfg.Redis)
	redisCtx, redisCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer redisCancel()
	if err := redisClient.Ping(redisCtx); err != nil {
		logger.Errorf(err, "failed to connect to redis")
		os.Exit(1)
	}
	cacheRepo := redis.NewCacheRepo(redisClient, infoConv, cfg.Redis, logger)

	conn, err := grpc.NewClient(
		cfg.Ml.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // явное указание gRPC-клиенту использовать НЕзащищённое соединение (без TLS).
	)
	if err != nil {
		logger.Errorf(err, "failed to initialize grpc client")
		os.Exit(1)
	}
	defer conn.Close()

	mlClient := proto.NewMachineLearningServiceClient(conn)
	ml := ml_service.NewMLService(mlClient, cfg.Ml, logger)
	imagesInfra := minioInfra.NewMinioInfrastructure(imageRepo, cfg.Minio, logger)

	producer, err := kafka.NewProducer(logger)
	if err != nil {
		logger.Errorf(err, "failed to initialize kafka producer")
		os.Exit(1)
	}

	productUC := usecase.NewProductUC(
		productRepo,
		categoryRepo,
		db.Pool,
		ml,
		imagesInfra,
		embRepo,
		logger,
		cacheRepo,
		prEmbeddingVersionRepo,
		producer,
	)

	grpcSrv := v1Grpc.NewGRPCServer(cfg.Grpc)
	grpcSrv.RegisterServices(productUC, logger)

	grpcErrCh := make(chan error, 1)
	go func() {
		logger.Infof("gRPC server starting on %s:%s", cfg.Grpc.NetworkMode, cfg.Grpc.Port)
		if err := grpcSrv.Start(); err != nil {
			logger.Errorf(err, "gRPC server failed")
			grpcErrCh <- err
		}
	}()

	r := chi.NewRouter()
	router := v1Http.NewRouter(r, logger)
	router.Init(productUC)

	httpSrv := v1Http.NewServer(r, cfg.Http)

	errCh := make(chan error, 1)
	go func() {
		logger.Infof("HTTP server started on port %s", cfg.Http.Port)
		if err := httpSrv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf(err, "HTTP server failed: %v", err)
			errCh <- err
		}
	}()

	// === Ожидание сигнала или ошибки ===
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	var appErr error
	select {
	case appErr = <-errCh:
		logger.Errorf(appErr, "HTTP server fatal error")
	case appErr = <-grpcErrCh:
		logger.Errorf(appErr, "gRPC server fatal error")
	case <-shutdown:
		logger.Infof("Received shutdown signal, stopping gracefully...")
	}

	// === Graceful shutdown ===
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpSrv.Stop(shutdownCtx); err != nil {
		logger.Errorf(err, "HTTP server shutdown error")
	} else {
		logger.Infof("HTTP server stopped")
	}

	if err := grpcSrv.Stop(shutdownCtx); err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			logger.Errorf(err, "gRPC server shutdown error")
		} else {
			logger.Warnf("gRPC server shutdown timeout")
		}
	} else {
		logger.Infof("gRPC server stopped")
	}

	done := make(chan error, 1)
	go func() {
		done <- imagesInfra.WaitForCleanup(shutdownCtx) // internal cleanup, без таймаута
	}()

	select {
	case err := <-done:
		if err != nil {
			logger.Warnf("MinIO cleanup error: %v", err)
		} else {
			logger.Infof("MinIO cleanup completed")
		}
	case <-time.After(5 * time.Second): // локальный таймаут ожидания cleanup
		logger.Warnf("MinIO cleanup did not finish before shutdown, some temporary objects may remain")
	}

	if qdrantClient != nil {
		if err := qdrantClient.Client.Close(); err != nil {
			logger.Warnf("Qdrant close error: %v", err)
		}
	}

	if redisClient != nil {
		if err := redisClient.Client.Close(); err != nil {
			logger.Warnf("Redis close error: %v", err)
		}
	}

	if db != nil {
		db.Close()
	}

	logger.Infof("Application shutdown complete")
	if appErr != nil {
		os.Exit(1)
	}
}

func initPGDB(logger logger.Logger, cfg *config.Config) (*postgres.PgDatabase, error) {
	db, err := postgres.Connect(cfg.Db)
	if err != nil {
		logger.Errorf(err, "failed to connect to database")
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	if err := db.RunMigrations(logger); err != nil {
		logger.Errorf(err, "failed to run migrations")
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	if err := db.Ping(); err != nil {
		logger.Errorf(err, "failed to ping database")
		return nil, e.Wrap(whereami.WhereAmI(), err)
	}

	return db, nil
}
