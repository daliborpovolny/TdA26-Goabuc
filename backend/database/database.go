package database

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	database "tourbackend/database/gen"
	"tourbackend/utils"

	"database/sql"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

func Initialize() (*sql.DB, *database.Queries) {

	resetDB := os.Getenv("RESET_DB")

	if resetDB == "true" {
		os.Remove("./database/db_file.db")
	}

	ctx := context.Background()

	db, err := sql.Open("sqlite", "./database/db_file.db")
	if err != nil {
		panic(err)
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}

	queries := database.New(db)

	if resetDB == "true" {
		Seed(queries)
		fmt.Println("db reseted and seeded")
	}

	return db, queries
}

// this inserts initial data into database which is useful for testing
func Seed(queries *database.Queries) {

	ctx := context.Background()

	hash, err := utils.HashPassword("12345678")
	if err != nil {
		panic(err)
	}

	_, err = queries.CreateUser(ctx, database.CreateUserParams{
		FirstName: "Adminov",
		LastName:  "Adminsky",
		Email:     "adminov.adminksky@goabuc.cz",
		Hash:      hash,
	})
	if err != nil {
		panic(err)
	}
}
