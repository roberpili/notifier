package seed

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func randomUsername() string {
	return fmt.Sprintf("test_user_%d", rand.IntN(10000))
}

func SeedUsers(ctx context.Context, tx pgx.Tx) (int64, error) {
	query := `
		INSERT INTO users (
			username, email, password_hash, role, is_active, created_at
 		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	 `

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var userID int64
	username := randomUsername()
	err = tx.QueryRow(
		ctx,
		query,
		username,
		username+"@gmail.com",
		string(hashedPassword),
		"user",
		true,
		time.Now(),
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("inserting seed user failed: %w", err)
	}

	return userID, nil
}
