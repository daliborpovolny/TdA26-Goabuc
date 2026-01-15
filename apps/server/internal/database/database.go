package database

import (
	"context"
	"database/sql"
	_ "embed"
	"os"

	gen "tourbackend/internal/database/gen"

	_ "modernc.org/sqlite"
)

var PATH_TO_DB string = "../../internal/database/db_file.db"

//go:embed schema.sql
var ddl string

func Initialize(resetDB bool) (*sql.DB, *gen.Queries) {

	PATH_TO_DB_ENV := os.Getenv("PATH_TO_DB")
	if PATH_TO_DB_ENV != "" {
		PATH_TO_DB = PATH_TO_DB_ENV
	}

	if resetDB {
		os.Remove(PATH_TO_DB)
	}

	ctx := context.Background()

	db, err := sql.Open("sqlite", PATH_TO_DB)
	if err != nil {
		panic(err)
	}

	if resetDB {
		// create tables
		if _, err := db.ExecContext(ctx, ddl); err != nil {
			panic(err)
		}
	}

	queries := gen.New(db)
	return db, queries
}
