package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"

	"train_http_storage/internals/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.postgres.NewStorage"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	absPath, err := filepath.Abs("internals/storages/postgres/init.sql")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	initSQl, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(string(initSQl))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) UploadFiles(file models.FileToAdd) error {
	const op = "storage.UploadFiles"

	stmt, err := s.db.Prepare(`INSERT INTO files (
                  alias, path_to_file, user_id) VALUES ($1,$2,$3)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(file.Alias, file.PathToFile, file.UserId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) NewUser(name string) error {
	const op = "storage.postgres.NewUser"

	stmt, err := s.db.Prepare(`INSERT INTO users (name) values ($1)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) TakeAliasesByUserId(id int64) ([]string, error) {
	const op = "storage.postgres.TakeAliasesByUserId"

	stmt, err := s.db.Prepare(`SELECT alias from files where user_id=$1`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %s", err.Error())
		}
	}()

	aliases := make([]string, 0)

	for rows.Next() {
		var alias string
		if err = rows.Scan(&alias); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		aliases = append(aliases, alias)
	}

	return aliases, nil
}

func (s *Storage) CloseStorage() error {
	const op = "storage.CloseStorage"

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// todo unload ( из бд) логика сжатия в хенделере

// todo integralTestForDb bd + write tests
