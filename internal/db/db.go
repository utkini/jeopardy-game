package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB(storagePath string) (*sqlx.DB, error) {
	const op = "storage.sqlite.New"

	db, err := sqlx.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS players(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		score INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE questions (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    category TEXT NOT NULL,
	    question TEXT NOT NULL,
	    answer TEXT NOT NULL,
	    points INTEGER NOT NULL
    );`
	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return db, nil
}
