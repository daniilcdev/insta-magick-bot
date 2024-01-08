package telegram

import "github.com/daniilcdev/insta-magick-bot/generated/queries"

type Storage interface {
	NewRequest(file, requesterId string)

	Schedule(limit int64) []string
	GetCompleted() []queries.ObtainCompletedRow
	RemoveCompleted()

	CompleteRequests(files []string)
}
