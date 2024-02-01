package telegram

import (
	"context"
	"fmt"
	"path"

	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ReplyToPhotoHandler struct {
	scheduler WorkScheduler
	storage   Storage
	log       Logger
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

func (h *ReplyToPhotoHandler) WithScheduler(scheduler WorkScheduler) *ReplyToPhotoHandler {
	h.scheduler = scheduler
	return h
}

func (h *ReplyToPhotoHandler) WithStorage(storage Storage) *ReplyToPhotoHandler {
	h.storage = storage
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

	filterName := getFilterNameFromTextEntities(update.Message)
	filter, err := h.storage.FindFilter(filterName)

	if err != nil {
		h.log.ErrStr(fmt.Sprintf("failed to find filter %s", filterName))
		bot.SendMessage(ctx,
			&tg.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text: `Упс!
Указан недоступный фильтр.

Чтобы узнать доступные фильтры,
используйте команду /filters`,
			},
		)
		return
	}

	go bot.SendMessage(ctx,
		&tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: `Отлично!
Вскоре я отправлю Вам результат обработки.`,
		},
	)

	fileId, err := getFileId(msg)
	if err != nil {
		h.log.Err(err)

		bot.SendMessage(ctx,
			&tg.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text: `Упс!
Проблема с получением файла.
Попробуйте отправить файл заново.`,
			},
		)
		return
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
	h.scheduler.Schedule(types.Work{
		File:        file.FileID + path.Ext(dlLink),
		RequesterId: fmt.Sprintf("%d", update.Message.Chat.ID),
		Filter:      types.Instructions(filter.Receipt),
		URL:         dlLink,
	})
}
