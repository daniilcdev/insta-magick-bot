package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/daniilcdev/insta-magick-bot/adapters"
	"github.com/daniilcdev/insta-magick-bot/client/telegram"
	"github.com/daniilcdev/insta-magick-bot/config"
	"github.com/daniilcdev/insta-magick-bot/internal"
	messaging "github.com/daniilcdev/insta-magick-bot/messaging/pkg"
	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
)

func main() {
	mq := messaging.InitMessageQueue()
	defer mq.Close()

	workDone := make(chan *types.Work)
	mq.Notify(internal.WorkDone, workDone)
	defer close(workDone)

	cfg := config.LoadConfig()
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
		WithWorkScheduler(mq)

	sendBackAdapter := &adapters.SendFileBackHandler{
		Log:        adapters.NewLogger().WithTag("SendbackAdapter"),
		Client:     botClient,
		Storage:    db,
		ResultsDir: cfg.ResultsDir(),
	}
	go sendBackAdapter.ListenResult(ctx, workDone)

	go botClient.Start(db)

	waitForExit()
}

func waitForExit() {
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	<-interupt
}
