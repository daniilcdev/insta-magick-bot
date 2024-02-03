package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	botToken          string
	processedFilesDir string
	dbDriver          string
	dbConn            string
}

func LoadConfig() *AppConfig {
	godotenv.Load(
		"./env/private/telegram.env",
		"./env/private/db.env",
	)
	return &AppConfig{
		botToken:          os.Getenv("TELEGRAM_BOT_TOKEN"),
		processedFilesDir: os.Getenv("PROCESSED_FILES_DIR"),

		dbDriver: os.Getenv("DB_DRIVER"),
		dbConn:   os.Getenv("DB_CONN"),
	}
}

func (cfg *AppConfig) BotToken() string {
	return cfg.botToken
}

func (cfg *AppConfig) DbDriver() string {
	return cfg.dbDriver
}

func (cfg *AppConfig) DbConn() string {
	return cfg.dbConn
}

func (cfg *AppConfig) ResultsDir() string {
	return cfg.processedFilesDir
}