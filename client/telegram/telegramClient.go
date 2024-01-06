package telegram

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	folderscanner "github.com/daniilcdev/insta-magick-bot/workers/folderScanner"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramClient struct {
	log Logger
	bot *tg.Bot
	ctx context.Context
}

func NewBotClient(ctx context.Context, botToken string) (*TelegramClient, error) {
	if botToken == "" {
		panic(errors.New("missing token"))
	}

	tgc := TelegramClient{
		log: nil,
	}

	b, err := tg.New(botToken)
	if err != nil {
		return nil, err
	}

	b.RegisterHandler(tg.HandlerTypeMessageText, "/start", tg.MatchTypeExact, tgc.StartHandler)
	b.RegisterHandlerMatchFunc(tgc.PhotoMessageMatch, tgc.PhotoMessageHandler)

	tgc.bot = b
	tgc.ctx = ctx
	return &tgc, nil
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
			Text:   "Processing file, please wait...",
		},
	)

	go downloadPhoto(dlLink, filename)
}

func downloadPhoto(url, name string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non-200 response code")
	}

	file, err := os.Create("./res/raw/" + name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TelegramClient) StartHandler(ctx context.Context, bot *tg.Bot, update *models.Update) {
	bot.SendMessage(ctx, &tg.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `Отправьте изображение для обработки.
		В силу текущих ограничений, пожалуйста, отправляйте по одному изображению за раз.`,
	})
}

func (tc *TelegramClient) Start() {
	scanner_sendback := folderscanner.FileScanner{}
	scanner_sendback.FoundFileHandler = tc
	go scanner_sendback.KeepScanning(tc.ctx, os.Getenv("IM_OUT_DIR"), 3*time.Second)

	tc.log.Info("Bot started")
	tc.bot.Start(tc.ctx)
}
