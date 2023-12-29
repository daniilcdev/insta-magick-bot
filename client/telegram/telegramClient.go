package telegram

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramClient struct {
	log log.Logger
	bot *tg.Bot
	ctx context.Context
}

func NewClassroomTrackerBot(botToken string) (*TelegramClient, error) {
	tgc := TelegramClient{
		log: *log.Default(),
	}
	ctx := context.Background()
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

	bot.SendMessage(ctx,
		&tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Processing file, please wait...",
		},
	)

	go downloadPhoto(bot.FileDownloadLink(file), file.FileID)
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
	file, err := os.Create("res/raw/" + name + path.Ext(url))
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	log.Printf("file saved - %s\n", name)
	return nil
}

func (tc *TelegramClient) StartHandler(ctx context.Context, bot *tg.Bot, update *models.Update) {
	tc.log.Println(update.Message.From.ID)
	tc.log.Println(update.Message.From.Username)
	tc.log.Println(update.Message.From.LanguageCode)
}

func (tc *TelegramClient) Start() {
	tc.log.Println("[INFO] Bot started")
	tc.bot.Start(tc.ctx)
}
