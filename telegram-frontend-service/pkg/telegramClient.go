package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	logging "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramClient struct {
	log       logging.Logger
	bot       *bot.Bot
	scheduler WorkScheduler
	cfg       LaunchOptions
}

func NewBotClient(cfg LaunchOptions) *TelegramClient {
	tgc := TelegramClient{
		cfg: cfg,
	}

	return &tgc
}

func (tc *TelegramClient) WithToken(token string) *TelegramClient {
	if token == "" {
		panic(errors.New("missing token"))
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			jsonData, _ := json.Marshal(update)
			fmt.Println(string(jsonData))
		}),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return tc
	}

	tc.bot = b

	return tc
}

func (tc *TelegramClient) WithWorkScheduler(scheduler WorkScheduler) *TelegramClient {
	tc.scheduler = scheduler
	return tc
}

func (tc *TelegramClient) WithLogger(logger logging.Logger) *TelegramClient {
	tc.log = logger
	return tc
}

func (tc *TelegramClient) Start(ctx context.Context, handlers ...CommandHandler) {
	if tc.bot == nil {
		panic("can't start client - bot wasn't set")
	}

	tc.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, tc.startHandler)

	for _, h := range handlers {
		tc.bot.RegisterHandlerMatchFunc(h.WillHandle, h.Handle)
	}

	tc.log.Info("Bot started")
	tc.bot.Start(ctx)
}

func (tc *TelegramClient) startHandler(ctx context.Context, bot *bot.Bot, update *models.Update) {
	params := chatResponseParams(update)
	params.Text = `Добро пожаловать в Imbot!

Для начала, отправьте в этот чат Ваше фото и через несколько секунд я отвечу Вам обработанной версией.

Чтобы узнать доступные фильтры,
используйте команду /filters`

	bot.SendMessage(ctx, params)
}

func (tc *TelegramClient) sendPhoto(ctx context.Context, chatId string, inputFile models.InputFile) {
	params := &bot.SendPhotoParams{
		ChatID: chatId,
		Photo:  inputFile,
	}

	if _, err := tc.bot.SendPhoto(ctx, params); err != nil {
		tc.log.ErrStr(fmt.Sprintf("failed to send image back '%v'", err))
	}
}

func chatResponseParams(update *models.Update) *bot.SendMessageParams {
	return &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
	}
}
