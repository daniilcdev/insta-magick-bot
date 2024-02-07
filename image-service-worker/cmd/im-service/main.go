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

	cfg := config.Load()
	fsOK := directoryReachable(cfg.InDir()) &&
		directoryReachable(cfg.OutDir()) && directoryReachable(cfg.TempDir())

	if !fsOK {
		log.Fatalln("invalid file storage setup")
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

func directoryReachable(dir string) bool {
	_, err := os.Stat(dir)

	switch {
	case os.IsNotExist(err):
		log.Printf("%s - %v\n", dir, err)
		return false
	case err != nil:
		log.Printf("%s - %v\n", dir, err)
		return false
	default:
		return true
	}
}

func waitForInterrupt() {
	interrup := make(chan os.Signal, 1)
	signal.Notify(interrup, os.Interrupt, syscall.SIGTERM)

	s := <-interrup
	log.Printf("syscall: '%v'\n", s)
}
