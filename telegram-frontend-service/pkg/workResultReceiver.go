package telegram

import (
	"context"
	"fmt"
	"os"

	messaging "github.com/daniilcdev/insta-magick-bot/messaging/pkg"
	logging "github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/logger"
	"github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/internal/storage"
	"github.com/go-telegram/bot/models"
)

type WorkResultReceiver struct {
	tc      *TelegramClient
	storage storage.Storage
	log     logging.Logger
}

func NewResultReceiver(tc *TelegramClient) *WorkResultReceiver {
	return &WorkResultReceiver{
		tc: tc,
	}
}

func (receiver *WorkResultReceiver) WithLogger(logger logging.Logger) *WorkResultReceiver {
	receiver.log = logger
	return receiver
}
func (receiver *WorkResultReceiver) WithStorage(storage storage.Storage) *WorkResultReceiver {
	receiver.storage = storage
	return receiver
}

func (receiver *WorkResultReceiver) HandleCompletedWork(ctx context.Context, work *messaging.Work) {
	select {
	case <-ctx.Done():
		return
	default:
		if err := receiver.completeWorkRequest(ctx, work); err != nil {
			receiver.log.Err(err)
		}
	}
}

func (receiver *WorkResultReceiver) completeWorkRequest(ctx context.Context, work *messaging.Work) error {
	if work == nil {
		return fmt.Errorf("work is nil")
	}

	filePath := receiver.tc.cfg.ProcessedFilesDir + work.File
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	request, err := receiver.storage.CompleteRequest(work.RequestId)
	if err != nil {
		return err
	}

	if err = receiver.tc.sendPhoto(ctx,
		request.RequesterID,
		&models.InputFileUpload{
			Filename: work.File,
			Data:     file,
		}); err != nil {
		return err
	}

	defer receiver.storage.RemoveRequest(request.ID)

	file.Close()

	if err := os.Remove(filePath); err != nil {
		receiver.log.Warn(fmt.Sprintf("can't remove file: '%v'", err))
	}

	return err
}
