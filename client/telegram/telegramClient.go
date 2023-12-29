package telegram

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	folderscanner "github.com/daniilcdev/insta-magick-bot/workers/folderScanner"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramClient struct {
	log log.Logger
	bot *tg.Bot
	ctx context.Context
}

func NewBotClient(ctx context.Context, botToken string) (*TelegramClient, error) {
	if botToken == "" {
		panic(errors.New("missing token"))
	}

	tgc := TelegramClient{
		log: *log.Default(),
	}

	opts := []tg.Option{
		// bot.WithDefaultHandler(ctb.DefaultHandler),
	}

	b, err := tg.New(botToken, opts...)
	if err != nil {
		tgc.log.Println(err)
		return nil, err
	}

	b.RegisterHandler(tg.HandlerTypeMessageText, "/start", tg.MatchTypeExact, tgc.StartHandler)
	b.RegisterHandlerMatchFunc(tgc.PhotoMessageMatch, tgc.PhotoMessageHandler)

	tgc.bot = b
	tgc.ctx = ctx
	return &tgc, nil
}

func (tc *TelegramClient) PhotoMessageMatch(update *models.Update) bool {
	return len(update.Message.Photo) > 0
}

func (tc *TelegramClient) PhotoMessageHandler(ctx context.Context, bot *tg.Bot, update *models.Update) {
	defer mu.Unlock()

	params := tg.GetFileParams{}
	params.FileID = update.Message.Photo[3].FileID
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

	//Create a empty file
	file, err := os.Create("./res/raw/" + name)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	log.Printf("new raw file: %s\n", name)
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

	tc.log.Println("[INFO] Bot started")
	tc.bot.Start(tc.ctx)
}
