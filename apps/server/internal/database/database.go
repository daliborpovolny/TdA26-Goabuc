package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"time"

	gen "tourbackend/internal/database/gen"
	"tourbackend/internal/utils"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

var PATH_TO_DB string = "../../internal/database/db_file.db"

//go:embed schema.sql
var ddl string

func Initialize() (*sql.DB, *gen.Queries) {

	PATH_TO_DB_ENV := os.Getenv("PATH_TO_DB")
	if PATH_TO_DB_ENV != "" {
		PATH_TO_DB = PATH_TO_DB_ENV
	}

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

	hash, err := utils.HashPassword("1234")
	if err != nil {
		panic(err)
	}

	_, err = queries.CreateUser(ctx, gen.CreateUserParams{
		FirstName: "Adminov",
		LastName:  "Adminsky",
		Email:     "ad.min@goabuc.cz",
		Hash:      hash,
	})
	if err != nil {
		panic(err)
	}

	now := time.Now().Unix()
	_, err = queries.CreateCourse(ctx, gen.CreateCourseParams{
		Uuid:        uuid.NewString(),
		Name:        "Pottery for Beginners",
		Description: "Intro into the wonderful world of pottery. No matter you experience you are welcome!",
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	now = time.Now().Unix()
	_, err = queries.CreateCourse(ctx, gen.CreateCourseParams{
		Uuid:        uuid.NewString(),
		Name:        "Potions 101",
		Description: "Intro into potion making, fast-paced course for serious sorcerers only",
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	now = time.Now().Unix()
	_, err = queries.CreateCourse(ctx, gen.CreateCourseParams{
		Uuid:        uuid.NewString(),
		Name:        "Zebra Riding Advanced",
		Description: "A guide to advanced zebra riding techniques, must already own a zebra",
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}
