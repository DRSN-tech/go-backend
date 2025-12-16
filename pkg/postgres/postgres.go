package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/DRSN-tech/go-backend/pkg/e"
	"github.com/DRSN-tech/go-backend/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// PgDatabase инкапсулирует подключение к PostgreSQL и управление миграциями.
type PgDatabase struct {
	Pool *pgxpool.Pool
	Dsn  string
}

func NewPgDatabase(pool *pgxpool.Pool, dsn string) *PgDatabase {
	return &PgDatabase{Pool: pool, Dsn: dsn}
}

// Connect устанавливает соединение с PostgreSQL с использованием переменных окружения.
func Connect() (*PgDatabase, error) {
	const op = "PgDatabase.Connect"
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("SSL_MODE"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewPgDatabase(pool, dsn), nil
}

// Close корректно закрывает пул соединений к базе данных.
func (db *PgDatabase) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// RunMigrations применяет ожидающие миграции из директории db/migrations.
func (db *PgDatabase) RunMigrations(logger logger.Logger) error {
	const (
		op                 = "PgDatabase.RunMigrations"
		driverName         = "pgx"
		databaseDriverName = "postgres"
		sourceURL          = "file://db/migrations"
	)

	sqlDb, err := sql.Open(driverName, db.Dsn)
	if err != nil {
		return err
	}
	defer sqlDb.Close()

	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		databaseDriverName,
		driver,
	)
	if err != nil {
		return e.Wrap(op, err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return e.Wrap(op, err)
	}

	logger.Infof("migrations applied successfully")
	return nil
}
