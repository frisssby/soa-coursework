package db

import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func NewDatabase(uri string) (*sql.DB, error) {
	db, err := sql.Open("clickhouse", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to clickhouse: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to check connection: %w", err)
	}
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}
	return db, nil
}

func createTables(db *sql.DB) error {
	var err error
	if _, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS views (
			user_id String,
			task_id String,
			author_id String,
			PRIMARY KEY (user_id, task_id)
		) engine=ReplacingMergeTree`,
	); err != nil {
		return err
	}
	if _, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS likes (
			user_id String,
			task_id String,
			author_id String,
			PRIMARY KEY (user_id, task_id)
		) engine=ReplacingMergeTree`,
	); err != nil {
		return err
	}
	return nil
}
