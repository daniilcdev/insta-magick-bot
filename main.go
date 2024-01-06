package main

import (
	"context"
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
	godotenv.Load("./env/private/telegram.env", "./env/imagemagick.env")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner_receive := folderscanner.FileScanner{}
	scanner_receive.FoundFileHandler = imclient.NewProcessor(os.Getenv("IM_IN_DIR"), os.Getenv("IM_OUT_DIR"))
	go scanner_receive.KeepScanning(ctx, "./res/raw/", 20*time.Second)

	botClient := telegram.NewBotClient(ctx)

	scanner_sendback := folderscanner.FileScanner{}
	scanner_sendback.FoundFileHandler = &adapters.SendFileBackHandler{
		Log:    &adapters.DefaultLoggerAdapter{},
		Client: botClient,
	}
	go scanner_sendback.KeepScanning(ctx, os.Getenv("IM_OUT_DIR"), 30*time.Second)

	go botClient.
		WithToken(os.Getenv("TELEGRAM_BOT_TOKEN")).
		WithLogger(&adapters.DefaultLoggerAdapter{}).
		Start()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	<-interupt
}
