package telegram

import (
	"errors"
	"log"
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

func getFilterNameFromTextEntities(m *models.Message) string {
	if m == nil {
		log.Default().Println("getFilterNameFromTextEntities - message is nil")
		return ""
	}

	for i := 0; i < len(m.Entities); i++ {
		ent := m.Entities[i]
		switch ent.Type {
		case models.MessageEntityTypeBotCommand:
			stripOff := ent.Length + ent.Offset + 1
			if stripOff >= len(m.Text) {
				return ""
			}

			return m.Text[ent.Length+ent.Offset+1:]
		default:
			continue
		}
	}

	return ""
}
