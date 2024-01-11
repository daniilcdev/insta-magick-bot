package adapters

import "github.com/daniilcdev/insta-magick-bot/client/telegram"

type LogBuilder interface {
	telegram.Logger
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
