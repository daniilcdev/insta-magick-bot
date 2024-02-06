package telegram

import (
	"context"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type photoMessageHandler struct {
}

func NewPhotoMessageHandler() *photoMessageHandler {
	return &photoMessageHandler{}
}

func (tc *photoMessageHandler) WillHandle(update *models.Update) bool {
	return tc.photoMessageMatch(update)
}

func (tc *photoMessageHandler) Handle(ctx context.Context, bot *tg.Bot, update *models.Update) {
	tc.handlePhotoMessage(ctx, bot, update)
}

func (tc *photoMessageHandler) photoMessageMatch(update *models.Update) bool {
	return hasPhotoAttached(update.Message)
}

func (tc *photoMessageHandler) handlePhotoMessage(ctx context.Context, bot *tg.Bot, update *models.Update) {
	bot.SendMessage(ctx,
		&tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: `Отлично!
Теперь ответьте на свою фотографию,
добавив команду /filter и название фильтра, например:
/filter Dark Indie`,
		},
	)
}
