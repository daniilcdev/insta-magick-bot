package config

import (
	"os"

	imclient "github.com/daniilcdev/insta-magick-bot/client/imClient"
	"github.com/daniilcdev/insta-magick-bot/client/telegram"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	imclient.IMConfig
	telegram.BotConfig

	inDir    string
	outDir   string
	botToken string
	dbDriver string
	dbConn   string
}

func LoadConfig() *AppConfig {
	godotenv.Load(
		"./env/private/telegram.env",
		"./env/private/db.env",
		"./env/imagemagick.env",
	)
	return &AppConfig{
		botToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		inDir:    os.Getenv("IM_IN_DIR"),
		outDir:   os.Getenv("IM_OUT_DIR"),
		dbDriver: os.Getenv("DB_DRIVER"),
		dbConn:   os.Getenv("DB_CONN"),
	}
}

func (cfg *AppConfig) BotToken() string {
	return cfg.botToken
}

func (cfg *AppConfig) DownloadDir() string {
	return cfg.inDir
}

func (cfg *AppConfig) InDir() string {
	return cfg.inDir
}

func (cfg *AppConfig) OutDir() string {
	return cfg.outDir
}

func (cfg *AppConfig) DbDriver() string {
	return cfg.dbDriver
}

func (cfg *AppConfig) DbConn() string {
	return cfg.dbConn
}