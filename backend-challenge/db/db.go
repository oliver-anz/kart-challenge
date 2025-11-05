package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{sqlDB}

	if err := db.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

func (db *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		category TEXT NOT NULL,
		price REAL NOT NULL,
		image_thumbnail TEXT,
		image_mobile TEXT,
		image_tablet TEXT,
		image_desktop TEXT
	);

	CREATE TABLE IF NOT EXISTS valid_coupons (
		code TEXT PRIMARY KEY,
		occurrence_count INTEGER NOT NULL
	);
	`

	_, err := db.Exec(schema)
	return err
}
