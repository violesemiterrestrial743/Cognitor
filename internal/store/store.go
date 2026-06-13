package store

import (
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`pragma foreign_keys = on`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := Migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func encode(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func decode[T any](raw string) (T, error) {
	var value T
	if raw == "" {
		raw = "null"
	}
	err := json.Unmarshal([]byte(raw), &value)
	return value, err
}
