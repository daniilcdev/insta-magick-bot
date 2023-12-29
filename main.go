package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

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
	go scanner_receive.KeepScanning(ctx, "./res/raw/", 2*time.Second)

	botClient, err := telegram.NewClassroomTrackerBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Default().Println(err)
	}

	go botClient.Start()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	<-interupt
}
