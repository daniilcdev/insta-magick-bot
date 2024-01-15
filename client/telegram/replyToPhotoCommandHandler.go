package telegram

import (
	"context"
	"fmt"
	"path"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ReplyToPhotoHandler struct {
	imageLoader *imageWebLoader
	log         Logger
}

func NewReplyToPhoto() *ReplyToPhotoHandler {
	return &ReplyToPhotoHandler{
		log: &nopLoggerAdapter{},
	}
}

func (h *ReplyToPhotoHandler) WithLogger(log Logger) *ReplyToPhotoHandler {
	h.log = log
	return h
}

func (h *ReplyToPhotoHandler) WithImageLoader(loader *imageWebLoader) *ReplyToPhotoHandler {
	h.imageLoader = loader
	return h
}

func (h *ReplyToPhotoHandler) WillHandle(update *models.Update) bool {
	msg := update.Message
	switch {
	case msg == nil:
		return false
	case hasPhotoAttached(msg.ReplyToMessage):
		return true
	default:
		return false
	}
}

func (h *ReplyToPhotoHandler) Handle(ctx context.Context, bot *tg.Bot, update *models.Update) {
	msg := update.Message.ReplyToMessage

	fileId, err := getFileId(msg)
	if err != nil {
		h.log.Err(err)
		return
	}

	go bot.SendMessage(ctx,
		&tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: `Отлично!
Вскоре я отправлю Вам результат обработки.`,
		},
	)

	var cmd string

	for i := 0; i < len(update.Message.Entities) && cmd == ""; i++ {
		ent := update.Message.Entities[i]
		switch ent.Type {
		case models.MessageEntityTypeBotCommand:
			cmd = update.Message.Text[ent.Length+ent.Offset+1:]
		default:
			continue
		}
	}

	params := tg.GetFileParams{}
	params.FileID = fileId
	file, err := bot.GetFile(ctx, &params)
	if err != nil {
		bot.SendMessage(ctx, &tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Oops!\n" + err.Error(),
		},
		)
		return
	}

	dlLink := bot.FileDownloadLink(file)
	dlParams := downloadParams{
		url:         dlLink,
		outFilename: file.FileID + path.Ext(dlLink),
		requesterId: fmt.Sprintf("%d", update.Message.Chat.ID),
		filter:      cmd,
	}

	go h.imageLoader.downloadPhoto(dlParams)
}
