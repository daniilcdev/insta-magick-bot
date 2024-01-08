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

func (s *SqliteStorage) GetPendingRequests(limit int64) []string {
	rows, err := s.q.ObtainPendingFiles(context.Background(), limit)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		return nil
	}

	return rows
}

func (s *SqliteStorage) GetRequester(file string) (string, error) {
	result, err := s.q.GetRequest(context.Background(), file)
	if err != nil {
		return "", err
	}

	return result.RequesterID, nil
}

func (s *SqliteStorage) GetRequestersByFilenames(files []string) ([]queries.GetRequestersByFilenamesRow, error) {
	result, err := s.q.GetRequestersByFilenames(context.Background(), files)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SqliteStorage) UpdateFilesStatus(files []string) {
	err := s.q.UpdateFilesStatus(context.Background(), files)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}
}

func (s *SqliteStorage) RemoveRequest(file string) {
	err := s.q.DeleteRequest(context.Background(), file)
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
