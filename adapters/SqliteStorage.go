package adapters

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/daniilcdev/insta-magick-bot/generated/queries"
	_ "github.com/mattn/go-sqlite3"
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

func OpenStorageConnection() (*SqliteStorage, error) {
	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_CONN"))
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

func (s *SqliteStorage) NewRequest(file, requesterId string) {
	err := s.q.CreateRequest(
		context.Background(),
		queries.CreateRequestParams{
			File:        file,
			RequesterID: requesterId,
		},
	)

	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}
}

func (s *SqliteStorage) Schedule(limit int64) []string {
	rows, err := s.q.SchedulePending(context.Background(), limit)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		return nil
	}

	return rows
}

func (s *SqliteStorage) GetCompleted() []queries.GetRequestsInStatusRow {
	rows, err := s.q.GetRequestsInStatus(context.Background(), string(Completed))
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		return nil
	}

	return rows
}

func (s *SqliteStorage) RemoveCompleted() {
	err := s.q.DeleteRequestsInStatus(context.Background(), string(Completed))
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}
}

func (s *SqliteStorage) CompleteRequests(files []string) {
	args := queries.UpdateRequestsStatusParams{
		Filenames: files,
		Status:    string(Completed),
	}
	err := s.q.UpdateRequestsStatus(context.Background(), args)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}
}

func (s *SqliteStorage) Close() error {
	defer func() {
		log.Println("db connection closed")
	}()

	return s.db.Close()
}
