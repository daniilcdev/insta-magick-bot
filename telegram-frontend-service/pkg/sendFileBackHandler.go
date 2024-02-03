package telegram

import (
	"os"

	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
	"github.com/go-telegram/bot/models"
)

func (sb *TelegramClient) ListenResult(result chan *types.Work) {
	for {
		select {
		case work := <-result:
			sb.do(work)
		case <-sb.ctx.Done():
			return
		}
	}
}

func (sb *TelegramClient) do(work *types.Work) {
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

	sb.sendPhoto(work.RequesterId, &models.InputFileUpload{
		Filename: work.RequesterId,
		Data:     f,
	})
}
