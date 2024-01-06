package telegram

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramClient struct {
	log       Logger
	bot       *tg.Bot
	imgLoader *imageWebLoader
	ctx       context.Context
}

func NewBotClient(ctx context.Context) *TelegramClient {
	tgc := TelegramClient{
		log: nil,
	}

	tgc.ctx = ctx
	tgc.log = &nopLoggerAdapter{}
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

	b.RegisterHandler(tg.HandlerTypeMessageText, "/start", tg.MatchTypeExact, tc.StartHandler)
	b.RegisterHandlerMatchFunc(tc.PhotoMessageMatch, tc.PhotoMessageHandler)

	tc.bot = b

	return tc
}

func (tc *TelegramClient) WithLogger(logger Logger) *TelegramClient {
	tc.log = logger
	return tc
}

func (tc *TelegramClient) PhotoMessageMatch(update *models.Update) bool {
	switch {
	case update.Message.Document != nil && strings.HasPrefix(update.Message.Document.MimeType, "image"):
		return true
	case len(update.Message.Photo) > 0:
		return true
	default:
		return false
	}
}

func (tc *TelegramClient) PhotoMessageHandler(ctx context.Context, bot *tg.Bot, update *models.Update) {
	params := tg.GetFileParams{}
	fileId, err := getFileId(update.Message)
	if err != nil {
		tc.log.Err(err.Error())
		return
	}

	params.FileID = fileId
	file, err := bot.GetFile(ctx, &params)
	if err != nil {
		bot.SendMessage(ctx, &tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed to download image: " + err.Error(),
		},
		)
	}

	dlLink := bot.FileDownloadLink(file)
	filename := file.FileID + path.Ext(dlLink)

	defer mu.Unlock()
	mu.Lock()
	imgToChatMap[filename] = update.Message.Chat.ID

	bot.SendMessage(ctx,
		&tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Result will be sent back shortly...",
		},
	)

	dlParams := downloadParams{
		url:         dlLink,
		outFilename: file.FileID + path.Ext(dlLink),
		outDir:      "./res/raw/",
		requesterId: fmt.Sprintf("%d", update.Message.Chat.ID),
	}

	go tc.imgLoader.downloadPhoto(dlParams)
}

func (tc *TelegramClient) StartHandler(ctx context.Context, bot *tg.Bot, update *models.Update) {
	bot.SendMessage(ctx, &tg.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `Отправьте изображение для обработки.
		В силу текущих ограничений, пожалуйста, отправляйте по одному изображению за раз.`,
	})
}

func (tc *TelegramClient) Start() {
	if tc.bot == nil {
		panic("can't start client - bot wasn't set")
	}

	tc.log.Info("Bot started")
	tc.bot.Start(tc.ctx)
}
