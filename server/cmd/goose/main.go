package main

import (
	"database/sql"
	"embed"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	// setup database
	dburl := os.Getenv("DB_URL")
	if dburl == "" {
		panic("DB_URL is empty")
	}
	db, err := sql.Open("sqlite3", dburl)

	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	// run app
}

