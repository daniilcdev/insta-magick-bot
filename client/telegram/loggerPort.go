package telegram

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Err(msg string)
}
