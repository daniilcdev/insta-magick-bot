package telegram

type BotConfig interface {
	BotToken() string
	ResultsDir() string
}
