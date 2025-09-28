package storage

import (
	"context"
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"

	"lucasbonna/pulse/db"
)

//go:embed schema.sql
var ddl string

func NewSQLiteDB() (*db.Queries, error) {
	ctx := context.Background()

	startedDb, err := sql.Open("sqlite", "db.sqlite?_journal=WAL&_timeout=5000&_synchronous=NORMAL")
	if err != nil {
		return nil, err
	}

	startedDb.SetMaxOpenConns(1)
	startedDb.SetMaxIdleConns(1)
	startedDb.SetConnMaxLifetime(0)

	if _, err := startedDb.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := db.New(startedDb)

	return queries, nil
}
