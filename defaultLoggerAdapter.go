package main

import "log"

type DefaultLoggerAdapter struct {
}

func (logger *DefaultLoggerAdapter) Err(msg string) {
	log.Default().Printf("[ERROR] %s\n", msg)
}

func (logger *DefaultLoggerAdapter) Warn(msg string) {
	log.Default().Printf("[WARN] %s\n", msg)
}

func (logger *DefaultLoggerAdapter) Info(msg string) {
	log.Default().Printf("[INFO] %s\n", msg)
}
