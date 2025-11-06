package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	// Add busy_timeout to handle concurrent access
	sqlDB, err := sql.Open("sqlite3", dbPath+"?_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for SQLite
	// SQLite serializes writes, so high connection counts provide no benefit
	// and can cause "database is locked" errors under contention
	// For read-heavy workloads: 3-5 connections is optimal
	// For write-heavy workloads: 1-2 connections is better
	sqlDB.SetMaxOpenConns(5) // Reduced from 25 - SQLite doesn't benefit from high concurrency
	sqlDB.SetMaxIdleConns(2) // Keep 2 connections warm for quick reuse
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(30 * time.Second) // Close idle connections after 30s

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{sqlDB}, nil
}
