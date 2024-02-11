package telegram

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/config"
	logging "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramClient struct {
	log logging.Logger
	bot *bot.Bot
	cfg *config.AppConfig
}

func NewBotClient(cfg *config.AppConfig) (*TelegramClient, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			jsonData, _ := json.Marshal(update)
			fmt.Println(string(jsonData))
		}),
	}

	botApi, err := bot.New(cfg.BotToken, opts...)
	if err != nil {
		return nil, err
	}

	tgc := TelegramClient{
		cfg: cfg,
		bot: botApi,
	}

	return &tgc, nil
}

func (tc *TelegramClient) WithLogger(logger logging.Logger) *TelegramClient {
	tc.log = logger
	return tc
}

func (tc *TelegramClient) WithCommandHandler(handler CommandHandler) *TelegramClient {
	tc.bot.RegisterHandlerMatchFunc(handler.WillHandle, handler.Handle)
	return tc
}

func (tc *TelegramClient) Start(ctx context.Context) error {
	if tc.bot == nil {
		return fmt.Errorf("can't start client - bot wasn't set")
	}

	tc.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, tc.startHandler)

	tc.log.Info("Bot started")
	go tc.bot.Start(ctx)
	return nil
}

func (tc *TelegramClient) startHandler(ctx context.Context, bot *bot.Bot, update *models.Update) {
	params := ChatResponseParams(update)
	params.Text = `Добро пожаловать в Imbot!

Для начала, отправьте в этот чат Ваше фото и через несколько секунд я отвечу Вам обработанной версией.

Чтобы узнать доступные фильтры,
используйте команду /filters`

	bot.SendMessage(ctx, params)
}

func (tc *TelegramClient) sendPhoto(ctx context.Context, chatId string, inputFile models.InputFile) error {
	params := &bot.SendPhotoParams{
		ChatID: chatId,
		Photo:  inputFile,
	}

	if _, err := tc.bot.SendPhoto(ctx, params); err != nil {
		return err
	}

	return nil
}

func ChatResponseParams(update *models.Update) *bot.SendMessageParams {
	return &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
	}
}
