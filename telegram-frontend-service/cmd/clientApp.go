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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var mq messaging.MessagingClient
	var err error
	if mq, err = messaging.Connect(ctx); err != nil {
		return nil
	}

	defer mq.Close()

	db, err := storage.Connect(app.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	logger := logging.NewLogger().WithTag("BotClient")

	botClient, err := telegram.NewBotClient(app.cfg)
	if err != nil {
		return err
	}

	botClient.WithLogger(logger).
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

	if err = botClient.Start(ctx); err != nil {
		return err
	}

	receiver := telegram.NewResultReceiver(botClient).
		WithLogger(logging.NewLogger().WithTag("WorkDone Listener")).
		WithStorage(db)

	mq.OnMessage(messaging.WorkDone, receiver.HandleCompletedWork)

	waitForExit()
	return nil
}

func waitForExit() {
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	<-interupt
}
