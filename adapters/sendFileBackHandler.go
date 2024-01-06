package adapters

import (
	"fmt"
	"io/fs"
	"os"
	"sync"

	"github.com/daniilcdev/insta-magick-bot/client/telegram"
	"github.com/go-telegram/bot/models"
)

type SendFileBackHandler struct {
	Log    telegram.Logger
	Client *telegram.TelegramClient
}

func (sb *SendFileBackHandler) ProcessNewFilesInDir(dir string, entries []fs.DirEntry) {
	telegram.Mu.Lock()

	wg := sync.WaitGroup{}
	for _, entry := range entries {
		fileName := entry.Name()

		chatId, ok := telegram.ImgToChatMap[fileName]

		if !ok {
			sb.Log.Warn(fmt.Sprintf("chatId not found for file %s", fileName))
			continue
		}

		f, err := os.Open(dir + fileName)
		if err != nil {
			sb.Log.Err(fmt.Sprintf("can't open file %v", err))
			continue
		}
		defer f.Close()

		wg.Add(1)
		go sb.Client.SendPhoto(&wg, chatId, &models.InputFileUpload{
			Filename: fileName,
			Data:     f,
		})

		delete(telegram.ImgToChatMap, fileName)
		os.Remove(dir + fileName)
	}

	telegram.Mu.Unlock()

	wg.Wait()
}
