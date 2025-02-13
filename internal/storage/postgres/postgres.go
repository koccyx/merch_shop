package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/koccyx/avito_assignment/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrWrongId = errors.New("wrong id")
	ErrEmptyFields = errors.New("empty fields")
	ErrNotFound = errors.New("url not found")
	ErrEntryExists = errors.New("entry exists")
)

func prepareDb(db *sql.DB) error{

    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id VARCHAR PRIMARY KEY,
            username VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
			balance INTEGER DEFAULT 1000
        );`)
    if err != nil {
        return fmt.Errorf("error executing SQL query: %w", err)
    }

	return nil
}

func New(cfg *config.Config) (*sql.DB, error){	
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", cfg.Storage.User, cfg.Storage.Password, cfg.Storage.Addres, cfg.Storage.Port, cfg.Storage.Database, cfg.Storage.Schema)
	db, err := sql.Open("pgx", connStr)

	if err != nil {
		return nil, fmt.Errorf("%s: opening db error: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: ping error: %w", op, err)
	}

	err = prepareDb(db)

	if err != nil {
		return nil, fmt.Errorf("%s: error during db opening: %w", op, err)
	}

	return db, nil
}