package telegram

import (
	"context"
	"strings"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (tc *TelegramClient) matchListFiltersCommand(update *models.Update) bool {
	switch {
	case update.Message == nil:
		return false
	case strings.HasPrefix(update.Message.Text, "/filters"):
		return true
	default:
		return false
	}
}

func (tc *TelegramClient) handleListFiltersCommand(ctx context.Context, bot *tg.Bot, update *models.Update) {
	_, err := bot.SendMessage(ctx, &tg.SendMessageParams{
		Text:   strings.Join(tc.filtersPool, "\n"),
		ChatID: update.Message.Chat.ID,
	})

	if err != nil {
		tc.log.Err(err.Error())
	}
}
