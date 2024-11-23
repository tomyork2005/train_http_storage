package postgreSQL

import (
	"RESTarchive/internals/storages"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.PostgreSQL.NewStorage"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "op", err)
	}

	// todo add index to FK
	stmt1, err := db.Prepare(`CREATE TABLE IF NOT EXISTS users ( 
    id SERIAL PRIMARY KEY, 
    name VARCHAR(100) UNIQUE NOT NULL);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt2, err := db.Prepare(`CREATE TABLE IF NOT EXISTS files ( 
    id SERIAL PRIMARY KEY, 
    alias VARCHAR(100) UNIQUE, 
    path_to_file VARCHAR(200) UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id));`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt1.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt2.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) UploadFiles(file storages.FileToAdd) error {
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
	const op = "storage.NewUser"

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

// todo unload ( из бд) логика сжатия в хенделере

// todo test bd + write tests
