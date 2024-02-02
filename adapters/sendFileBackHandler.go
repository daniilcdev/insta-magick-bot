package adapters

import (
	"context"
	"os"

	"github.com/daniilcdev/insta-magick-bot/client/telegram"
	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
	"github.com/go-telegram/bot/models"
)

type SendFileBackHandler struct {
	Log        telegram.Logger
	Client     *telegram.TelegramClient
	Storage    telegram.Storage
	ResultsDir string
}

func (sb *SendFileBackHandler) ListenResult(ctx context.Context, result chan *types.Work) {
	for {
		select {
		case work := <-result:
			sb.do(work)
		case <-ctx.Done():
			return
		}
	}
}

func (sb *SendFileBackHandler) do(work *types.Work) {
	filePath := sb.ResultsDir + work.File
	f, err := os.Open(filePath)
	if err != nil {
		sb.Log.Err(err)
		return
	}
	defer f.Close()

	sb.Client.SendPhoto(work.RequesterId, &models.InputFileUpload{
		Filename: work.RequesterId,
		Data:     f,
	})
}
