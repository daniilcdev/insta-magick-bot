package telegram

import "github.com/daniilcdev/insta-magick-bot/generated/queries"

type Storage interface {
	NewRequest(file, requesterId string)
	GetRequester(file string) (string, error)
	RemoveRequest(file string)
	GetPendingRequests(limit int64) []string
	UpdateFilesStatus(files []string)
	GetRequestersByFilenames(files []string) ([]queries.GetRequestersByFilenamesRow, error)
}
