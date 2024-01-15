package telegram

type BotConfig interface {
	BotToken() string
	DownloadDir() string
}
