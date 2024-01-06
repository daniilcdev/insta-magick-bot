package telegram

import (
	"errors"

	"github.com/go-telegram/bot/models"
)

func getFileId(m *models.Message) (string, error) {

	switch {
	case len(m.Photo) > 0:
		return m.Photo[len(m.Photo)-1].FileID, nil
	case m.Document != nil:
		return m.Document.FileID, nil
	}

	return "", errors.New("no file in message")
}
