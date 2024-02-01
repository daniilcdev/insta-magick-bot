package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/daniilcdev/insta-magick-bot/adapters"
	imclient "github.com/daniilcdev/insta-magick-bot/client/imClient"
	"github.com/daniilcdev/insta-magick-bot/client/telegram"
	"github.com/daniilcdev/insta-magick-bot/config"
	"github.com/daniilcdev/insta-magick-bot/workers"
)

func main() {
	cfg := config.LoadConfig()
	InitMessageQueue()

	db, err := adapters.OpenStorageConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	botClient := telegram.NewBotClient(ctx, cfg).
		WithToken(cfg.BotToken()).
		WithFiltersPool(db.FilterNames()).
		WithLogger(adapters.NewLogger().WithTag("BotClient")).
		WithStorage(db).
		WithWorkScheduler(mq)

	sendBackAdapter := &adapters.SendFileBackHandler{
		Log:     adapters.NewLogger().WithTag("SendbackAdapter"),
		Client:  botClient,
		Storage: db,
	}
	go botClient.Start()

	// clean up 'stale' completed images
	sendBackAdapter.OnProcessCompleted(cfg.OutDir())

	imc := imclient.NewProcessor(cfg, db).
		WithCompletionHandler(sendBackAdapter)

	scanner_receive := workers.PipelineTrigger{
		Handler: imc,
	}
	go scanner_receive.KeepScanning(ctx, cfg.InDir(), 30*time.Second)

	waitForExit()
}

func waitForExit() {
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	<-interupt
}
