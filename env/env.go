package env

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Fatalf(err.Error())
		log.Fatalf("Error loading .env file")
	}
}
