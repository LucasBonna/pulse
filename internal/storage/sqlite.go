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

	startedDb, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		return nil, err
	}

	if _, err := startedDb.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := db.New(startedDb)

	return queries, nil
}
