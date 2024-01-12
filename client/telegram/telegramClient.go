package telegram

import (
	"context"
	"errors"
	"fmt"
	"sync"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramClient struct {
	log         Logger
	bot         *tg.Bot
	imgLoader   *imageWebLoader
	ctx         context.Context
	filtersPool []string
	cfg         BotConfig
}

func NewBotClient(ctx context.Context, cfg BotConfig) *TelegramClient {
	tgc := TelegramClient{
		log: &nopLoggerAdapter{},
		cfg: cfg,
	}

	tgc.ctx = ctx
	tgc.imgLoader = &imageWebLoader{}
	return &tgc
}

func (tc *TelegramClient) WithToken(token string) *TelegramClient {
	if token == "" {
		panic(errors.New("missing token"))
	}

	b, err := tg.New(token)
	if err != nil {
		return tc
	}

	tc.bot = b

	return tc
}

func (tc *TelegramClient) WithLogger(logger Logger) *TelegramClient {
	tc.log = logger
	return tc
}

func (tc *TelegramClient) WithStorage(storage Storage) *TelegramClient {
	tc.imgLoader.storage = storage
	return tc
}

func (tc *TelegramClient) Start() {
	if tc.bot == nil {
		panic("can't start client - bot wasn't set")
	}

	tc.bot.RegisterHandler(tg.HandlerTypeMessageText, "/start", tg.MatchTypeExact, tc.startHandler)
	tc.bot.RegisterHandlerMatchFunc(tc.matchListFiltersCommand, tc.handleListFiltersCommand)
	tc.bot.RegisterHandlerMatchFunc(tc.photoMessageMatch, tc.photoMessageHandler)

	tc.log.Info("Bot started")
	tc.bot.Start(tc.ctx)
}

func (tc *TelegramClient) WithFiltersPool(pool []string) *TelegramClient {
	tc.filtersPool = pool
	return tc
}

func (tc *TelegramClient) SendPhoto(wg *sync.WaitGroup, chatId any, inputFile models.InputFile) {
	defer wg.Done()

	params := &tg.SendPhotoParams{
		ChatID: chatId,
		Photo:  inputFile,
	}

	_, err := tc.bot.SendPhoto(tc.ctx, params)

	if err != nil {
		tc.log.ErrStr(fmt.Sprintf("failed to send image back %v", err))
		return
	}
}

func (tc *TelegramClient) startHandler(ctx context.Context, bot *tg.Bot, update *models.Update) {
	bot.SendMessage(ctx, &tg.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `Отправьте изображение для обработки.
		В силу текущих ограничений, пожалуйста, отправляйте по одному изображению за раз.`,
	})
}
