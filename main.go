package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	folderscanner.FoundFileHandler = imclient.NewProcessor(os.Getenv("IM_IN_DIR"), os.Getenv("IM_OUT_DIR"))
	go folderscanner.KeepScanning(ctx, "./res/raw", 2*time.Second)

	botClient, err := telegram.NewClassroomTrackerBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Default().Println(err)
	}

	go botClient.Start()

	fmt.Println("keeping system alive for 10 minutes")
	<-time.After(10 * time.Second)
}
