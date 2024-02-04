package telegram

import (
	"context"
	"os"

	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
	"github.com/go-telegram/bot/models"
)

func (sb *TelegramClient) ListenResult(ctx context.Context, result chan *types.Work) {
	for {
		select {
		case work := <-result:
			sb.do(ctx, work)
		case <-ctx.Done():
			return
		}
	}
}

func (sb *TelegramClient) do(ctx context.Context, work *types.Work) {
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
