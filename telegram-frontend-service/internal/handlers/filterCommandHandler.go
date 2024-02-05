package telegram

import (
	"context"
	"strings"

	logging "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/logger"
	telegram "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/pkg"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type listFiltersHandler struct {
	storage telegram.Storage
	log     logging.Logger
}

func NewListFiltersHandler() *listFiltersHandler {
	return &listFiltersHandler{}
}

func (h *listFiltersHandler) WithStorage(storage telegram.Storage) *listFiltersHandler {
	h.storage = storage
	return h
}
func (h *listFiltersHandler) WithLogger(logger logging.Logger) *listFiltersHandler {
	h.log = logger
	return h
}

func (h *listFiltersHandler) WillHandle(update *models.Update) bool {
	return matchListFiltersCommand(update)
}

func (h *listFiltersHandler) Handle(ctx context.Context, bot *tg.Bot, update *models.Update) {
	h.handleListFiltersCommand(ctx, bot, update)
}

func matchListFiltersCommand(update *models.Update) bool {
	switch {
	case update.Message == nil:
		return false
	case strings.HasPrefix(update.Message.Text, "/filters"):
		return true
	default:
		return false
	}
}

func (h *listFiltersHandler) handleListFiltersCommand(ctx context.Context, bot *tg.Bot, update *models.Update) {
	availableFilters := h.storage.FilterNames()
	msgParams := &tg.SendMessageParams{
		ChatID: update.Message.Chat.ID,
	}

	if len(availableFilters) == 0 {
		msgParams.Text = "Отсутствуют доступные фильтры"
	} else {
		msgParams.Text = strings.Join(availableFilters, "\n")
	}

	_, err := bot.SendMessage(ctx, msgParams)

	if err != nil {
		h.log.Err(err)
	}
}