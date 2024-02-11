package telegram

import (
	"context"
	"strings"

	logging "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/logger"
	"github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/storage"
	telegram "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/pkg"
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type listFiltersHandler struct {
	storage storage.FiltersProvider
	log     logging.Logger
}

func NewListFiltersHandler() *listFiltersHandler {
	return &listFiltersHandler{}
}

func (h *listFiltersHandler) WithStorage(storage storage.FiltersProvider) *listFiltersHandler {
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
	msgParams := telegram.ChatResponseParams(update)

	if len(availableFilters) == 0 {
		msgParams.Text = "Отсутствуют доступные фильтры"
	} else {
		msgParams.Text = strings.Join(availableFilters, "\n")
	}

	if _, err := bot.SendMessage(ctx, msgParams); err != nil {
		h.log.Err(err)
	}
}
