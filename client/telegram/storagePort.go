package telegram

import "github.com/daniilcdev/insta-magick-bot/generated/queries"

type NewRequest struct {
	File        string
	RequesterId string
	Filter      string
}

type Storage interface {
	CreateRequest(newRequest *NewRequest)

	Schedule(limit int64) []queries.SchedulePendingRow
	GetCompleted() []queries.GetRequestsInStatusRow
	RemoveCompleted()

	CompleteRequests(files []string)

	FindFilter(name string) (filter queries.Filter, err error)
}
