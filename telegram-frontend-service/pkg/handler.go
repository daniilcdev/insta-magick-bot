package telegram

import (
	"context"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CommandHandler interface {
	WillHandle(update *models.Update) bool
	Handle(ctx context.Context, bot *tg.Bot, update *models.Update)
}
