package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	Port      string
	AppUrl    string
	PgUrl     string
	RedisUrl  string
	SecretKey string
}

var Env *Environment

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Env = &Environment{
		Port:      getEnv("PORT", "8000"),
		AppUrl:    getEnv("APP_URL", "localhost:8000"),
		PgUrl:     getEnv("PG_URL", "postgres://unbit:unbit@localhost:5432/unbit?sslmode=disable"),
		RedisUrl:  getEnv("REDIS_URL", "redis://localhost:6379"),
		SecretKey: getEnv("SECRET_KEY", "your-insecure-default-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
