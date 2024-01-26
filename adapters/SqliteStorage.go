package adapters

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/daniilcdev/insta-magick-bot/client/telegram"
	"github.com/daniilcdev/insta-magick-bot/config"
	"github.com/daniilcdev/insta-magick-bot/generated/queries"
	_ "github.com/lib/pq"
)

type SqliteStorage struct {
	db *sql.DB
	q  *queries.Queries
}

type requestStatus string

var (
	Pending    requestStatus = "Pending"
	InProgress requestStatus = "Processing"
	Completed  requestStatus = "Completed"
)

func OpenStorageConnection(cfg *config.AppConfig) (*SqliteStorage, error) {
	db, err := sql.Open(cfg.DbDriver(), cfg.DbConn())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("db connected, ping -> success")

	q := queries.New(db)

	return &SqliteStorage{db: db, q: q}, nil
}

func (s *SqliteStorage) FilterNames() []string {
	names, err := s.q.GetNames(context.Background())
	if err != nil {
		return []string{}
	}

	return names
}

func (s *SqliteStorage) CreateRequest(newRequest *telegram.NewRequest) {
	err := s.q.CreateRequest(
		context.Background(),
		queries.CreateRequestParams{
			File:        newRequest.File,
			RequesterID: newRequest.RequesterId,
			FilterName:  newRequest.Filter,
		},
	)

	reportErr(err)
}

func (s *SqliteStorage) Schedule(limit int64) []queries.SchedulePendingRow {
	rows, err := s.q.SchedulePending(context.Background(), limit)
	if err != nil {
		log.Printf("[ERROR] (Schedule) - '%v'\n", err)
		return nil
	}

	return rows
}

func (s *SqliteStorage) GetCompleted() []queries.GetRequestsInStatusRow {
	rows, err := s.q.GetRequestsInStatus(context.Background(), string(Completed))
	if err != nil {
		log.Printf("[ERROR] (GetCompleted) - '%v'\n", err)
		return nil
	}

	return rows
}

func (s *SqliteStorage) RemoveCompleted() {
	err := s.q.DeleteRequestsInStatus(context.Background(), string(Completed))
	reportErr(err)

}

func (s *SqliteStorage) CompleteRequests(files []string) {
	args := queries.UpdateRequestsStatusParams{
		Filenames: files,
		Status:    string(Completed),
	}
	err := s.q.UpdateRequestsStatus(context.Background(), args)
	reportErr(err)

}

func (s *SqliteStorage) FindFilter(name string) (filter queries.Filter, err error) {
	switch name {
	case "":
		err = errors.New("filter name is empty")
		reportErr(err)
		return queries.Filter{}, err
	default:
		filter, err = s.q.GetReceiptWithName(context.Background(), name)
		reportErr(err)
		return filter, err
	}
}

func (s *SqliteStorage) Rollback(files []string) {
	err := s.q.UpdateRequestsStatus(context.Background(), queries.UpdateRequestsStatusParams{
		Filenames: files,
		Status:    string(Pending),
	})
	reportErr(err)
}

func (s *SqliteStorage) Close() error {
	defer func() {
		log.Println("db connection closed")
	}()

	return s.db.Close()
}

func reportErr(err error) {
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}
}
