package postgreSQL

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) (*Storage, error) {
	const op = "storage.PostgreSQL.NewStorage"

	// todo create postgreSQLConfig to connect
	db, err := sql.Open("postgres")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "op", err)
	}

	stmt1, err := db.Prepare(`CREATE TABLE IF NOT EXISTS users ( 
    id SERIAL PRIMARY KEY, 
    name VARCHAR(100) UNIQUE NOT NULL);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt2, err := db.Prepare(`CREATE TABLE IF NOT EXIST files ( 
    id SERIAL PRIMARY KEY, 
    alias VARCHAR(100) UNIQUE, 
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

// todo upload (в бд) логика сжатия в хенделере

// todo unload ( из бд) логика сжатия в хенделере

// todo test bd + write tests