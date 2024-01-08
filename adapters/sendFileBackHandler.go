package adapters

import (
	"fmt"
	"os"
	"sync"

	"github.com/daniilcdev/insta-magick-bot/client/telegram"
	"github.com/go-telegram/bot/models"
)

type SendFileBackHandler struct {
	Log     telegram.Logger
	Client  *telegram.TelegramClient
	Storage telegram.Storage
}

func (sb *SendFileBackHandler) ProcessNewFilesInDir(dir string, files []string) {
	defer func(d string, f []string) {
		sb.Storage.UpdateFilesStatus(files)

		for _, r := range f {
			os.Remove(d + r)
		}

	}(dir, files)

	responses, _ := sb.Storage.GetRequestersByFilenames(files)

	wg := sync.WaitGroup{}
	for _, r := range responses {
		f, err := os.Open(dir + r.File)
		if err != nil {
			sb.Log.Err(fmt.Sprintf("can't open file %v", err))
			continue
		}
		defer f.Close()

		wg.Add(1)
		go sb.Client.SendPhoto(&wg, r.RequesterID, &models.InputFileUpload{
			Filename: r.File,
			Data:     f,
		})
	}

	wg.Wait()
}
