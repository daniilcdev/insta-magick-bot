package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	messaging "github.com/daniilcdev/insta-magick-bot/messaging/pkg"
	"github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/config"
	handlers "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/handlers"
	logging "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/logger"
	"github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/storage"
	telegram "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/pkg"
	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
)

func main() {
	mq := messaging.InitMessageQueue()
	defer mq.Close()

	workDone := make(chan *types.Work)
	mq.Notify(messaging.WorkDone, workDone)
	defer close(workDone)

	cfg := config.LoadConfig()
	db, err := storage.OpenStorageConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.NewLogger().WithTag("BotClient")
	botClient := telegram.NewBotClient(cfg).
		WithToken(cfg.BotToken()).
		WithLogger(logger).
		WithWorkScheduler(mq)

	go botClient.ListenResult(ctx, workDone)

	h := getHandlers(logger, db, mq)
	go botClient.Start(
		ctx,
		h...,
	)

	waitForExit()
}

func getHandlers(
	logger logging.Logger,
	db telegram.Storage,
	mq telegram.WorkScheduler,
) []telegram.CommandHandler {
	h := []telegram.CommandHandler{
		handlers.NewPhotoMessageHandler(),
		handlers.NewReplyToPhoto().
			WithLogger(logger).
			WithStorage(db).
			WithScheduler(mq),

		handlers.
			NewListFiltersHandler().
			WithStorage(db).
			WithLogger(logger),
	}

	return h
}

func waitForExit() {
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	<-interupt
}
