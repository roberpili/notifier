package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rober0xf/notifier/internal/domain/repository"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *pgxpool.Pool

func InitPostgres() (*pgxpool.Pool, error) {
	config := GetConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DB_USER,
		url.QueryEscape(config.DB_PASS),
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging: %w", err)
	}

	// migrations
	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("error closing db: %v", err)
		}
	}()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("error with driver: %w", err)
	}

	migrationsPath := "file://" + GetEnvOrFatal("MIGRATIONS_PATH")
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("error while init migrate: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("error migrating: %w", err)
	}

	fmt.Println("database initialized")
	return pool, nil
}

func StartTokenCleanJob(ctx context.Context, userRepo repository.UserRepository) {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				deleted, err := userRepo.DeleteOldTokens(ctx)
				if err != nil {
					log.Printf("token cleanup failed: %v", err)
					continue
				}
				log.Printf("cleaned up %d expired tokens", deleted)

			case <-ctx.Done():
				log.Println("cleaned up finished")
				return
			}
		}
	}()
}
