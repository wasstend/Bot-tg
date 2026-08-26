package token

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func TokenMust() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env file found:", err)
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	if token == "" {
		log.Fatal("token is not set")
	}

	return token
}
