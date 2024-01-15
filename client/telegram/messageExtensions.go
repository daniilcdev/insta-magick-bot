package telegram

import (
	"errors"
	"strings"

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

func hasPhotoAttached(m *models.Message) bool {
	switch {
	case m == nil:
		return false
	case m.Document != nil && strings.HasPrefix(m.Document.MimeType, "image"):
		return true
	case len(m.Photo) > 0:
		return true
	default:
		return false
	}
}
