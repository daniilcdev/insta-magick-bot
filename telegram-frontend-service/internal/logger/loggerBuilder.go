package logging

type Logger interface {
	Info(msg string)
	Warn(msg string)
	ErrStr(msg string)
	Err(err error)
}

type LogBuilder interface {
	Logger
	WithTag(tag string) LogBuilder
}

func NewLogger() LogBuilder {
	return &defaultLoggerAdapter{
		tag: "default",
	}
}

func (logger *defaultLoggerAdapter) WithTag(tag string) LogBuilder {
	logger.tag = tag
	return logger
}
