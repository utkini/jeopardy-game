package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB(storagePath string) (*sqlx.DB, error) {
	const op = "storage.sqlite.New"
	dirStoragePath := filepath.Dir(storagePath)
	if _, err := os.Stat(dirStoragePath); os.IsNotExist(err) {
		os.Mkdir(dirStoragePath, os.ModePerm)
	}

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

	CREATE TABLE IF NOT EXISTS categories(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER NOT NULL,
		question TEXT NOT NULL,
		answer TEXT,
		points INTEGER NOT NULL,
		media_type TEXT NOT NULL DEFAULT 'none',
		media_url TEXT,
		is_answered BOOLEAN NOT NULL DEFAULT 0,
		was_correct BOOLEAN,
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_questions_category ON questions(category_id);
	CREATE INDEX IF NOT EXISTS idx_questions_answered ON questions(is_answered);`

	// Drop existing tables if they exist
	_, err = db.Exec("DROP TABLE IF EXISTS questions")
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS categories")
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS players")
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	// Create tables with new schema
	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return db, nil
}
