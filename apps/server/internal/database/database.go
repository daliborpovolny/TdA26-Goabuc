package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	gen "tourbackend/internal/database/gen"
	"tourbackend/internal/utils"

	_ "modernc.org/sqlite"
)

var PATH_TO_DB string = "../../internal/database/db_file.db"

//go:embed schema.sql
var ddl string

func Initialize() (*sql.DB, *gen.Queries) {

	resetDB := os.Getenv("RESET_DB")

	if resetDB == "true" {
		os.Remove(PATH_TO_DB)
	}

	ctx := context.Background()

	db, err := sql.Open("sqlite", PATH_TO_DB)
	if err != nil {
		panic(err)
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}

	queries := gen.New(db)

	if resetDB == "true" {
		Seed(queries)
		fmt.Println("db reseted and seeded")
	}

	return db, queries
}

// this inserts initial data into database which is useful for testing
func Seed(queries *gen.Queries) {

	ctx := context.Background()

	hash, err := utils.HashPassword("12345678")
	if err != nil {
		panic(err)
	}

	_, err = queries.CreateUser(ctx, gen.CreateUserParams{
		FirstName: "Adminov",
		LastName:  "Adminsky",
		Email:     "adminov.adminksky@goabuc.cz",
		Hash:      hash,
	})
	if err != nil {
		panic(err)
	}
}
