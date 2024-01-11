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
	"github.com/daniilcdev/insta-magick-bot/workers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(
		"./env/private/telegram.env",
		"./env/private/db.env",
		"./env/imagemagick.env",
	)

	db, err := adapters.OpenStorageConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	botClient := telegram.NewBotClient(ctx).
		WithToken(os.Getenv("TELEGRAM_BOT_TOKEN")).
		WithFiltersPool(db.FilterNames()).
		WithLogger(adapters.NewLogger().WithTag("BotClient")).
		WithStorage(db)

	sendBackAdapter := &adapters.SendFileBackHandler{
		Log:     adapters.NewLogger().WithTag("SendbackAdapter"),
		Client:  botClient,
		Storage: db,
	}
	go botClient.Start()

	// clean up 'stale' completed images
	sendBackAdapter.OnProcessCompleted(os.Getenv("IM_OUT_DIR"))

	imc := imclient.NewProcessor(os.Getenv("IM_OUT_DIR"), db).
		WithCompletionHandler(sendBackAdapter)

	scanner_receive := workers.PipelineTrigger{
		Handler: imc,
	}
	go scanner_receive.KeepScanning(ctx, os.Getenv("IM_IN_DIR"), 30*time.Second)

	waitForExit()
}

func waitForExit() {
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	<-interupt
}
