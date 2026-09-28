package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rober0xf/notifier/cmd/api"
	_ "github.com/rober0xf/notifier/docs"
)

// @title           Notifier API
// @version         1.0
// @description     Notifier API
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file loaded: %v", err)
	}

	server, err := api.NewAPIServer(":8080")
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	log.Fatal(server.Run())
}
