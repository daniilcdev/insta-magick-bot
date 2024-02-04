package telegram

import (
	"github.com/daniilcdev/insta-magick-bot/generated/queries"
)

type Storage interface {
	FilterNames() []string
	CreateRequest(newRequest *NewRequest)

	Schedule(limit int64) []queries.SchedulePendingRow
	Rollback(files []string)

	GetCompleted() []queries.GetRequestsInStatusRow
	RemoveCompleted()

	CompleteRequests(files []string)

	FindFilter(name string) (filter queries.Filter, err error)
}
