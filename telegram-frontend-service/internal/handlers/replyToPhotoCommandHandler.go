package telegram

import (
	"context"
	"fmt"
	"path"

	messaging "github.com/daniilcdev/insta-magick-bot/messaging/pkg"
	logging "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/logger"
	"github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/storage"
	pkg "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/pkg"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type replyToPhotoHandler struct {
	scheduler pkg.WorkScheduler
	storage   storage.Storage
	log       logging.Logger
}

func NewReplyToPhoto() *replyToPhotoHandler {
	return &replyToPhotoHandler{}
}

func (h *replyToPhotoHandler) WithLogger(log logging.Logger) *replyToPhotoHandler {
	h.log = log
	return h
}

func (h *replyToPhotoHandler) WithScheduler(scheduler pkg.WorkScheduler) *replyToPhotoHandler {
	h.scheduler = scheduler
	return h
}

func (h *replyToPhotoHandler) WithStorage(storage storage.Storage) *replyToPhotoHandler {
	h.storage = storage
	return h
}

func (h *replyToPhotoHandler) WillHandle(update *models.Update) bool {
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

func (h *replyToPhotoHandler) Handle(ctx context.Context, bot *tg.Bot, update *models.Update) {
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
	if file, err := bot.GetFile(ctx, &params); err != nil {
		bot.SendMessage(ctx, &tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Oops!\n" + err.Error(),
		})
	} else {
		dlLink := bot.FileDownloadLink(file)
		fileName := file.FileID + path.Ext(dlLink)

		requestId, err := h.storage.CreateRequest(fileName,
			fmt.Sprintf("%d", update.Message.Chat.ID),
			filter.Name)

		if err != nil {
			h.log.Err(err)
			return
		}

		work := messaging.Work{
			RequestId:   requestId,
			File:        fileName,
			Filter:      filter.Name,
			Instruction: messaging.Instructions(filter.Receipt),
			URL:         dlLink,
		}

		h.scheduler.Schedule(work)
	}
}
