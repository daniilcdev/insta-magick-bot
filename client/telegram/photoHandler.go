package telegram

import (
	"context"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tc *TelegramClient) photoMessageMatch(update *models.Update) bool {
	return hasPhotoAttached(update.Message)
}

func (tc *TelegramClient) photoMessageHandler(ctx context.Context, bot *tg.Bot, update *models.Update) {
	go bot.SendMessage(ctx,
		&tg.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: `Отлично!
Теперь ответьте на свою фотографию,
добавив команду /filter и название фильтра, например:
/filter Dark Indie`,
		},
	)
}
