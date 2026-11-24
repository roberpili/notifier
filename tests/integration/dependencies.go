package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rober0xf/notifier/internal/database"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type TestPayment struct {
	Name       string
	Amount     float64
	Type       string
	Category   string
	Date       string
	Paid       bool
	Recurrent  bool
	DueDate    *string
	PaidAt     *string
	Frequency  *string
	ReceiptURL *string
}

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if err := loadRootEnv(); err != nil {
		log.Printf("warning: could not load env file: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		database.GetEnvOrFatal[string]("POSTGRES_USER_TEST"),
		url.QueryEscape(database.GetEnvOrFatal[string]("POSTGRES_PASSWORD_TEST")),
		database.GetEnvOrFatal[string]("POSTGRES_HOST_TEST"),
		database.GetEnvOrFatal[string]("POSTGRES_PORT_TEST"),
		database.GetEnvOrFatal[string]("POSTGRES_DB_TEST"),
	)

	db, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)

	err = db.Ping(context.Background())
	require.NoError(t, err)

	runMigrations(t, db)
	cleanDatabase(t, db)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func insertTestUser(ctx context.Context, db *pgxpool.Pool, email, username, password string) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return 0, err
	}

	var id int
	err = db.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, username, email, string(hash), true).Scan(&id)
	return id, err
}

func insertTestPayment(
	ctx context.Context,
	db *pgxpool.Pool,
	userId int,
	p TestPayment,
) (int, error) {
	var id int
	err := db.QueryRow(ctx, `
		INSERT INTO payments (user_id, name, amount, type, category, date, paid, recurrent, due_date, paid_at, frequency, receipt_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`,
		userId,
		p.Name,
		p.Amount,
		p.Type,
		p.Category,
		p.Date,
		p.Paid,
		p.Recurrent,
		p.DueDate,
		p.PaidAt,
		p.Frequency,
		p.ReceiptURL).Scan(&id)
	return id, err
}

func insertVerificationToken(t *testing.T, db *pgxpool.Pool, userID int) string {
	t.Helper()

	rawToken := uuid.New().String()
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	_, err := db.Exec(
		context.Background(),
		`INSERT INTO user_tokens (user_id, token_hash, purpose, expires_at)
		 VALUES ($1, $2, 'email_verification', NOW() + INTERVAL '24 hours')`,
		userID, tokenHash,
	)
	require.NoError(t, err)

	return rawToken
}

func runMigrations(t *testing.T, db *pgxpool.Pool) {
	ctx := context.Background()

	_, _ = db.Exec(ctx, `DROP SCHEMA public CASCADE;`)
	_, _ = db.Exec(ctx, `CREATE SCHEMA public;`)

	users, err := os.ReadFile("../../internal/infraestructure/persistance/postgres/sqlc/schemas/schema_users.sql")
	require.NoError(t, err)

	_, err = db.Exec(ctx, string(users))
	require.NoError(t, err)

	user_tokens, err := os.ReadFile("../../internal/infraestructure/persistance/postgres/sqlc/schemas/schema_user_tokens.sql")
	require.NoError(t, err)

	_, err = db.Exec(ctx, string(user_tokens))
	require.NoError(t, err)

	payments, err := os.ReadFile("../../internal/infraestructure/persistance/postgres/sqlc/schemas/schema_payments.sql")
	require.NoError(t, err)

	_, err = db.Exec(ctx, string(payments))
	require.NoError(t, err)
}

func cleanDatabase(t *testing.T, db *pgxpool.Pool) {
	ctx := context.Background()

	_, err := db.Exec(ctx, "TRUNCATE users CASCADE")
	require.NoError(t, err)

	_, err = db.Exec(ctx, "TRUNCATE payments CASCADE")
	require.NoError(t, err)
}

func loadRootEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			log.Printf("loading .env from: %s", envPath)
			return godotenv.Load(envPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return fmt.Errorf(".env not found in any parent directory")
		}
		dir = parent
	}
}
