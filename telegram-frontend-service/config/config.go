package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	BotToken          string
	ProcessedFilesDir string
	DbDriver          string
	DbConn            string
}

func LoadConfig() *AppConfig {
	godotenv.Load(
		"./env/private/telegram.env",
		"./env/private/db.env",
	)
	return &AppConfig{
		BotToken:          os.Getenv("TELEGRAM_BOT_TOKEN"),
		ProcessedFilesDir: os.Getenv("PROCESSED_FILES_DIR"),
		DbDriver:          os.Getenv("DB_DRIVER"),
		DbConn:            os.Getenv("DB_CONN"),
	}
}
