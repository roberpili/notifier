package database

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type PostgresConfig struct {
	DB_NAME string
	DB_USER string
	DB_PASS string
	DB_HOST string
	DB_PORT string
}

var JwtKey []byte

func GetConfig() PostgresConfig {
	var config PostgresConfig
	config.DB_NAME = GetEnvOrFatal[string]("POSTGRES_DB")
	config.DB_USER = GetEnvOrFatal[string]("POSTGRES_USER")
	config.DB_PASS = GetEnvOrFatal[string]("POSTGRES_PASSWORD")
	config.DB_HOST = GetEnvOrFatal[string]("POSTGRES_HOST")
	config.DB_PORT = GetEnvOrFatal[string]("POSTGRES_PORT")

	// validate that port is a number
	if _, err := strconv.Atoi(config.DB_PORT); err != nil {
		log.Fatalf("PORT env is not a number: %v", err)
	}

	jwt_secret := strings.TrimSpace(GetEnvOrFatal[string]("JWT_SECRET"))
	if jwt_secret == "" {
		log.Fatalf("JWT_SECRET env is empty")
	}
	JwtKey = []byte(jwt_secret)

	return config
}

func GetEnvOrFatal[T string | int | bool](key string) T {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("env variable %s empty", key)
	}

	var result T

	switch any(result).(type) {
	case string:
		return any(value).(T)
	case int:
		v, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("invalid int for env %s: %v", key, err)
		}
		return any(v).(T)
	case bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("invalid bool for env %s: %v", key, err)
		}
		return any(v).(T)
	}

	panic("unsupported type")
}
