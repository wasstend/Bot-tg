package token

import (
	"os"

	"github.com/joho/godotenv"
)

func MustToken() string {
	if err := godotenv.Load(); err != nil {
		panic("no .env file found:" + err.Error())
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	if token == "" {
		panic("token is not set")
	}

	return token
}
