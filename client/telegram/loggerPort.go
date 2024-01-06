package telegram

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Err(msg string)
}

type nopLoggerAdapter struct {
}

func (nop *nopLoggerAdapter) Info(msg string) {}
func (nop *nopLoggerAdapter) Warn(msg string) {}
func (nop *nopLoggerAdapter) Err(msg string)  {}
