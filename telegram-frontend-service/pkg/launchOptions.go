package telegram

type LaunchOptions interface {
	BotToken() string
	ResultsDir() string
}
