package storage

import (
	"github.com/daniilcdev/insta-magick-bot/generated/queries"
)

type Storage interface {
	FiltersProvider

	CreateRequest(fileName, requester, filterName string) (int64, error)
	RemoveRequest(requestId int64) error
	CompleteRequest(requestId int64) (queries.Request, error)

	GetCompleted(limit int) ([]queries.GetRequestsInStatusRow, error)
}

type FiltersProvider interface {
	FilterNames() []string
	FindFilter(name string) (filter queries.Filter, err error)
}
