package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/daniilcdev/insta-magick-bot/generated/queries"
	"github.com/daniilcdev/insta-magick-bot/telegram-frontend-service/config"
	_ "github.com/lib/pq"
)

type sqlStorage struct {
	db *sql.DB
	q  *queries.Queries
}

type requestStatus string

var (
	Pending   requestStatus = "Pending"
	Completed requestStatus = "Completed"
)

func Connect(cfg *config.AppConfig) (*sqlStorage, error) {
	db, err := sql.Open(cfg.DbDriver, cfg.DbConn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("db connected, ping -> success")

	q := queries.New(db)

	return &sqlStorage{db: db, q: q}, nil
}

func (s *sqlStorage) FilterNames() []string {
	names, err := s.q.GetNames(context.Background())
	if err != nil {
		log.Default().Printf("can't get names: '%v'\n", err)
		return []string{}
	}

	return names
}

func (s *sqlStorage) CreateRequest(fileName, requester, filterName string) (int64, error) {
	return s.q.CreateRequest(
		context.Background(),
		queries.CreateRequestParams{
			File:        fileName,
			RequesterID: requester,
			FilterName:  filterName,
		},
	)
}

func (s *sqlStorage) GetCompleted(limit int) ([]queries.GetRequestsInStatusRow, error) {
	return s.q.GetRequestsInStatus(context.Background(), string(Completed))
}

func (s *sqlStorage) CompleteRequest(requestId int64) (queries.Request, error) {
	args := queries.UpdateRequestStatusParams{
		ID:     requestId,
		Status: string(Completed),
	}

	return s.q.UpdateRequestStatus(context.Background(), args)
}

func (s *sqlStorage) RemoveRequest(requestId int64) error {
	return s.q.DeleteRequest(context.Background(), requestId)
}

func (s *sqlStorage) FindFilter(name string) (filter queries.Filter, err error) {
	if name == "" {
		return queries.Filter{}, errors.New("filter name is empty")
	}

	return s.q.GetReceiptWithName(context.Background(), name)
}

func (s *sqlStorage) Close() error {
	defer log.Println("db connection closed")

	return s.db.Close()
}
