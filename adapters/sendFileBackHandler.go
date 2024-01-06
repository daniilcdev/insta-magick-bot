package adapters

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/daniilcdev/insta-magick-bot/client/telegram"
	"github.com/go-telegram/bot/models"
)

type SendFileBackHandler struct {
	Log    telegram.Logger
	Client *telegram.TelegramClient
}

func (sb *SendFileBackHandler) ProcessNewFile(dir string, entry fs.DirEntry) {
	defer telegram.Mu.Unlock()
	fileName := entry.Name()

	telegram.Mu.Lock()
	chatId, ok := telegram.ImgToChatMap[fileName]

	if !ok {
		sb.Log.Warn(fmt.Sprintf("chatId not found for file %s", fileName))
		return
	}

	f, err := os.Open(dir + fileName)
	if err != nil {
		sb.Log.Err(fmt.Sprintf("can't open file %v", err))
		return
	}
	defer f.Close()

	sb.Client.SendPhoto(chatId, &models.InputFileUpload{
		Filename: fileName,
		Data:     f,
	})

	delete(telegram.ImgToChatMap, fileName)
	os.Remove(dir + fileName)
}
