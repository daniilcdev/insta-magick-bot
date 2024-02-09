package main

import (
	"context"
	"os"
	"os/signal"

	messaging "github.com/daniilcdev/insta-magick-bot/messaging/pkg"
	"github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/config"
	handlers "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/handlers"
	logging "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/logger"
	"github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/storage"
	telegram "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/pkg"
)

type clientApp struct {
	cfg *config.AppConfig
}

func createApp() *clientApp {
	cfg := config.LoadConfig()
	return &clientApp{
		cfg: cfg,
	}
}

func (app *clientApp) start() error {
	var mq messaging.MessagingClient
	var err error
	if mq, err = messaging.Connect(); err != nil {
		return nil
	}

	defer mq.Close()

	workDone := make(chan *messaging.Work)
	mq.Notify(messaging.WorkDone, workDone)
	defer close(workDone)

	db, err := storage.Connect(app.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.NewLogger().WithTag("BotClient")
	botClient := telegram.NewBotClient(app.cfg).
		WithToken(app.cfg.BotToken()).
		WithLogger(logger).
		WithWorkScheduler(mq)

	botClient.
		WithCommandHandler(handlers.NewPhotoMessageHandler()).
		WithCommandHandler(handlers.NewReplyToPhoto().
			WithLogger(logger).
			WithStorage(db).
			WithScheduler(mq),
		).
		WithCommandHandler(handlers.NewListFiltersHandler().
			WithStorage(db).
			WithLogger(logger),
		)

	go botClient.ListenResult(ctx, workDone)
	go botClient.Start(ctx)

	waitForExit()
	return nil
}

func waitForExit() {
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	<-interupt
}
