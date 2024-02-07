package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	mq "github.com/daniilcdev/insta-magick-bot/image-service-worker/cmd/im-service-mq"
	"github.com/daniilcdev/insta-magick-bot/image-service-worker/config"
)

func main() {
	log.Println("starting worker...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid config: %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	imc := NewProcessor(cfg)

	workReceiver := &mq.MQWorkReceiver{
		W: imc,
	}

	defer workReceiver.Close()
	workReceiver.StartReceiving()
	log.Default().Println("worker started...")

	waitForInterrupt()
	cancel()
	<-ctx.Done()
}

func waitForInterrupt() {
	interrup := make(chan os.Signal, 1)
	signal.Notify(interrup, os.Interrupt, syscall.SIGTERM)

	s := <-interrup
	log.Printf("syscall: '%v'\n", s)
}
