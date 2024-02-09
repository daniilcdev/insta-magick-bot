package telegram

import (
	"context"
	"fmt"
	"os"

	messaging "github.com/daniilcdev/insta-magick-bot/messaging/pkg"
	"github.com/go-telegram/bot/models"
)

func (sb *TelegramClient) HandleCompletedWork(ctx context.Context, work *messaging.Work) {
	select {
	case <-ctx.Done():
		return
	default:
		if err := sb.do(ctx, work); err != nil {
			sb.log.Err(err)
		}
	}
}

func (sb *TelegramClient) do(ctx context.Context, work *messaging.Work) error {
	if work == nil {
		return fmt.Errorf("work is nil")
	}

	filePath := sb.cfg.ResultsDir() + work.File
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	err = sb.sendPhoto(ctx, work.RequesterId, &models.InputFileUpload{
		Filename: work.RequesterId,
		Data:     file,
	})

	file.Close()

	if err := os.Remove(filePath); err != nil {
		sb.log.Warn(fmt.Sprintf("can't remove file: '%v'", err))
	}

	return err
}
