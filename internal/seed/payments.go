package seed

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rober0xf/notifier/internal/domain/entity"
)

func SeedPayments(ctx context.Context, tx pgx.Tx, userID int64) error {
	query := `
		INSERT INTO payments (
			user_id, name, amount, type, category, date, due_date, paid, paid_at, recurrent, frequency
 		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
 `

	now := time.Now()
	_, err := tx.Exec(
		ctx,
		query,
		userID,
		"NBA Pass",
		15.00,
		entity.TransactionTypeExpense,
		entity.CategoryTypeSports,
		now.Format("2006-01-02"),
		now.AddDate(0, 1, 0).Format("2006-01-02"),
		false,
		nil,
		true,
		entity.FrequencyTypeMonthly,
	)

	return err
}
