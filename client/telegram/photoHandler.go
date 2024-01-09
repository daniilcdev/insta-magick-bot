package telegram

import (
	"context"
	"fmt"
	"path"
	"strings"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tc *TelegramClient) photoMessageMatch(update *models.Update) bool {
	switch {
	case update.Message.Document != nil && strings.HasPrefix(update.Message.Document.MimeType, "image"):
		return true
	case len(update.Message.Photo) > 0:
		return true
	default:
		return false
	}
}

func (tc *TelegramClient) photoMessageHandler(ctx context.Context, bot *tg.Bot, update *models.Update) {
	fileId, err := getFileId(update.Message)
	if err != nil {
		tc.log.Err(err.Error())
		return
	}

	go bot.SendMessage(ctx,
		&tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Result will be sent back shortly...",
		},
	)

	cmd := update.Message.Caption

	if strings.HasPrefix(cmd, "/filter ") {
		cmd = strings.TrimSpace(strings.TrimPrefix(cmd, "/filter"))
	}

	params := tg.GetFileParams{}
	params.FileID = fileId
	file, err := bot.GetFile(ctx, &params)
	if err != nil {
		bot.SendMessage(ctx, &tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed to request image file: " + err.Error(),
		},
		)
		return
	}

	dlLink := bot.FileDownloadLink(file)
	dlParams := downloadParams{
		url:         dlLink,
		outFilename: file.FileID + path.Ext(dlLink),
		outDir:      "./res/pending/",
		requesterId: fmt.Sprintf("%d", update.Message.Chat.ID),
		filter:      cmd,
	}

	go tc.imgLoader.downloadPhoto(dlParams)
}
