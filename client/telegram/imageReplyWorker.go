package telegram

import (
	"io/fs"
	"os"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var mu sync.Mutex = sync.Mutex{}
var imgToChatMap map[string]any = make(map[string]any)

func (tc *TelegramClient) ProcessNewFile(dir string, entry fs.DirEntry) {
	defer mu.Unlock()
	fileName := entry.Name()

	mu.Lock()
	chatId, ok := imgToChatMap[fileName]

	if !ok {
		tc.log.Printf("[WARN] chatId not found for file %s\n", fileName)
		return
	}

	f, err := os.Open(dir + fileName)
	if err != nil {
		tc.log.Printf("[ERROR] can't open file %v\n", err)
		return
	}
	defer f.Close()

	params := &bot.SendPhotoParams{
		ChatID: chatId,
		Photo: &models.InputFileUpload{
			Filename: fileName,
			Data:     f,
		},
	}

	_, err = tc.bot.SendPhoto(tc.ctx, params)

	if err != nil {
		tc.log.Printf("[ERROR] failed to send image back %v\n", err)
		return
	}

	delete(imgToChatMap, fileName)
	os.Remove(dir + fileName)
}
