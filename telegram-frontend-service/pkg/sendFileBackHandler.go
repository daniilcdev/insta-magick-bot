package telegram

import (
	"context"
	"os"

	messaging "github.com/daniilcdev/insta-magick-bot/messaging/pkg"
	"github.com/go-telegram/bot/models"
)

func (sb *TelegramClient) ListenResult(ctx context.Context, result chan *messaging.Work) {
	for {
		select {
		case work := <-result:
			sb.do(ctx, work)
		case <-ctx.Done():
			return
		}
	}
}

func (sb *TelegramClient) do(ctx context.Context, work *messaging.Work) {
	if work == nil {
		return
	}

	filePath := sb.cfg.ResultsDir() + work.File
	f, err := os.Open(filePath)
	if err != nil {
		sb.log.Err(err)
		return
	}
	defer f.Close()

	sb.sendPhoto(ctx, work.RequesterId, &models.InputFileUpload{
		Filename: work.RequesterId,
		Data:     f,
	})
}
