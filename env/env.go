package env

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env.development")
	if err != nil {
		slog.Info("Error loading .env.development file, Running in Prod mode")
	}
}
