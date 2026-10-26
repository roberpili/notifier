package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rober0xf/notifier/internal/seed"
	"github.com/rober0xf/notifier/pkg/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("no .env file loaded: %v", err)
	}

	if os.Getenv("ENV") == "production" {
		log.Fatal("seeding is not possible on production")
	}

	db, err := database.InitPostgres()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	if err := seed.Run(db); err != nil {
		log.Fatal(err)
	}
}
