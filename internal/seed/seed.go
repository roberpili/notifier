package seed

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(db *pgxpool.Pool) error {
	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting seed transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userID, err := SeedUsers(ctx, tx)
	if err != nil {
		return fmt.Errorf("seeding users: %w", err)
	}

	if err := SeedPayments(ctx, tx, userID); err != nil {
		return fmt.Errorf("seeding payments: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting seed transaction: %w", err)
	}

	return nil
}
