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
	folderscanner "github.com/daniilcdev/insta-magick-bot/workers/folderScanner"
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

	scanner_receive := folderscanner.FolderScanner{}
	scanner_receive.FoundFilesHandler = imclient.NewProcessor(os.Getenv("IM_OUT_DIR"))
	go scanner_receive.KeepScanning(ctx, os.Getenv("IM_IN_DIR"), 20*time.Second)

	botClient := telegram.NewBotClient(ctx)

	scanner_sendback := folderscanner.FolderScanner{}
	scanner_sendback.FoundFilesHandler = &adapters.SendFileBackHandler{
		Log:    &adapters.DefaultLoggerAdapter{},
		Client: botClient,
	}
	go scanner_sendback.KeepScanning(ctx, os.Getenv("IM_OUT_DIR"), 30*time.Second)

	go botClient.
		WithToken(os.Getenv("TELEGRAM_BOT_TOKEN")).
		WithLogger(&adapters.DefaultLoggerAdapter{}).
		WithStorage(db).
		Start()

	waitForExit()
}

func waitForExit() {
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	<-interupt
}
