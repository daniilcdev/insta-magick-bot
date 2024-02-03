package telegram

type Logger interface {
	Info(msg string)
	Warn(msg string)
	ErrStr(msg string)
	Err(err error)
}

type nopLoggerAdapter struct {
}

func (nop *nopLoggerAdapter) Info(msg string)   {}
func (nop *nopLoggerAdapter) Warn(msg string)   {}
func (nop *nopLoggerAdapter) ErrStr(msg string) {}
func (nop *nopLoggerAdapter) Err(err error)     {}
