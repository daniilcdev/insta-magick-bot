package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	pkg "github.com/daniilcdev/insta-magick-bot/client/telegram/pkg"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramClient struct {
	log         Logger
	bot         *tg.Bot
	scheduler   pkg.WorkScheduler
	ctx         context.Context
	filtersPool []string
	cfg         pkg.BotConfig
}

func NewBotClient(ctx context.Context, cfg pkg.BotConfig) *TelegramClient {
	tgc := TelegramClient{
		log: &nopLoggerAdapter{},
		cfg: cfg,
	}

	tgc.ctx = ctx
	return &tgc
}

func (tc *TelegramClient) WithToken(token string) *TelegramClient {
	if token == "" {
		panic(errors.New("missing token"))
	}

	opts := []tg.Option{
		tg.WithDefaultHandler(func(ctx context.Context, bot *tg.Bot, update *models.Update) {
			jsonData, _ := json.Marshal(update)
			fmt.Println(string(jsonData))
		}),
	}

	b, err := tg.New(token, opts...)
	if err != nil {
		return tc
	}

	tc.bot = b

	return tc
}

func (tc *TelegramClient) WithWorkScheduler(scheduler pkg.WorkScheduler) *TelegramClient {
	tc.scheduler = scheduler
	return tc
}

func (tc *TelegramClient) WithLogger(logger Logger) *TelegramClient {
	tc.log = logger
	return tc
}

func (tc *TelegramClient) Start(storage Storage) {
	if tc.bot == nil {
		panic("can't start client - bot wasn't set")
	}

	tc.bot.RegisterHandler(tg.HandlerTypeMessageText, "/start", tg.MatchTypeExact, tc.startHandler)
	tc.bot.RegisterHandlerMatchFunc(tc.matchListFiltersCommand, tc.handleListFiltersCommand)
	tc.bot.RegisterHandlerMatchFunc(tc.photoMessageMatch, tc.photoMessageHandler)

	replyHandler := NewReplyToPhoto().
		WithLogger(tc.log).
		WithStorage(storage).
		WithScheduler(tc.scheduler)

	tc.bot.RegisterHandlerMatchFunc(replyHandler.WillHandle, replyHandler.Handle)

	tc.log.Info("Bot started")
	tc.bot.Start(tc.ctx)
}

func (tc *TelegramClient) WithFiltersPool(pool []string) *TelegramClient {
	tc.filtersPool = pool
	return tc
}

func (tc *TelegramClient) startHandler(ctx context.Context, bot *tg.Bot, update *models.Update) {
	bot.SendMessage(ctx, &tg.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `Добро пожаловать в Imbot!

Для начала, отправьте в этот чат Ваше фото и через несколько секунд я отвечу Вам обработанной версией.

Чтобы узнать доступные фильтры,
используйте команду /filters`,
	})
}

func (tc *TelegramClient) sendPhoto(chatId string, inputFile models.InputFile) {
	params := &tg.SendPhotoParams{
		ChatID: chatId,
		Photo:  inputFile,
	}

	if _, err := tc.bot.SendPhoto(tc.ctx, params); err != nil {
		tc.log.ErrStr(fmt.Sprintf("failed to send image back %v", err))
	}
}
