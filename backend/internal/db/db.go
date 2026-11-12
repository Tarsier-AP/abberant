package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"abberant/backend/internal/config"
)

//go:embed all:migrations
var migrationFS embed.FS

type Close func()

func Connect(cfg config.Config) (*pgxpool.Pool, Close, error) {
	dsn := cfg.DSN

	var stopPostgres func()

	if dsn == "" {
		postgres := embeddedpostgres.NewDatabase(
			embeddedpostgres.DefaultConfig().
				Version(embeddedpostgres.V16).
				Port(uint32(cfg.EmbeddedPort)).
				Database("abberant").
				Username("postgres").
				Password("postgres").
				RuntimePath(cfg.DataDir),
		)

		if err := postgres.Start(); err != nil {
			return nil, nil, fmt.Errorf("start embedded postgres: %w", err)
		}

		stopPostgres = func() {
			if err := postgres.Stop(); err != nil {
				log.Printf("embedded postgres stop: %v", err)
			}
		}

		dsn = (&url.URL{
			Scheme: "postgres",
			User:   url.UserPassword("postgres", "postgres"),
			Host:   fmt.Sprintf("localhost:%d", cfg.EmbeddedPort),
			Path:   "abberant",
		}).String() + "?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		stopEmbedded(stopPostgres)
		return nil, nil, fmt.Errorf("open pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := waitForDatabase(ctx, pool); err != nil {
		stopEmbedded(stopPostgres)
		pool.Close()
		return nil, nil, err
	}

	if err := migrateUp(pool); err != nil {
		stopEmbedded(stopPostgres)
		pool.Close()
		return nil, nil, err
	}

	return pool, func() {
		pool.Close()

		if stopPostgres != nil {
			stopPostgres()
		}
	}, nil
}

func waitForDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	for {
		err := pool.Ping(ctx)
		if err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return fmt.Errorf("ping database: %w", err)
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("ping database: %w", err)
		case <-time.After(250 * time.Millisecond):
		}
	}
}

func stopEmbedded(stop func()) {
	if stop != nil {
		stop()
	}
}

// migrateUp applies every embedded migration exactly once via golang-migrate.
func migrateUp(pool *pgxpool.Pool) error {
	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("read embedded migrations: %w", err)
	}

	connConfig := pool.Config().ConnConfig

	sqlDB := stdlib.OpenDB(*connConfig)
	defer sqlDB.Close()

	driver, err := migratepgx.WithInstance(sqlDB, &migratepgx.Config{})
	if err != nil {
		return fmt.Errorf("prepare migration driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "abberant", driver)
	if err != nil {
		return fmt.Errorf("prepare migrator: %w", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
