package adapters

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
}

var mu sync.Mutex = sync.Mutex{}
var imgToChatMap map[string]string = make(map[string]string)

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
	return &SqliteStorage{db: db}, nil
}

func (s *SqliteStorage) NewRequest(file, requesterId string) {
	defer mu.Unlock()
	mu.Lock()

	imgToChatMap[file] = requesterId
}

func (s *SqliteStorage) GetRequester(file string) (string, error) {
	defer mu.Unlock()
	mu.Lock()

	result, ok := imgToChatMap[file]
	if !ok {
		return "", fmt.Errorf("no file requester, file: %s", file)
	}

	return result, nil
}

func (s *SqliteStorage) RemoveRequest(file string) {
	defer mu.Unlock()
	mu.Lock()

	delete(imgToChatMap, file)
}

func (s *SqliteStorage) Close() error {
	defer func() {
		log.Println("db connection closed")
	}()

	return s.db.Close()
}
