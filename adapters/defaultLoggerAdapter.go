package adapters

import (
	"log"
)

type defaultLoggerAdapter struct {
	tag string
}

func (logger *defaultLoggerAdapter) Err(msg string) {
	log.Default().Printf("(%s) -error- %s\n", logger.tag, msg)
}

func (logger *defaultLoggerAdapter) Warn(msg string) {
	log.Default().Printf("(%s) -warn- %s\n", logger.tag, msg)
}

func (logger *defaultLoggerAdapter) Info(msg string) {
	log.Default().Printf("(%s) -info- %s\n", logger.tag, msg)
}
