package main

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migtationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage", "", "Path to the SQLite database file")
	flag.StringVar(&migtationsPath, "migrations", "", "Path to the migrations directory")
	flag.StringVar(&migrationsTable, "table", "migrations", "Name of the migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is required")
	}

	if migtationsPath == "" {
		panic("migrations path is required")
	}

	m, err := migrate.New(
		"file://"+migtationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable))

	if err != nil {
		panic(fmt.Sprintf("failed to create migrate instance: %v", err))
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No new migrations to apply")
		}

		panic(fmt.Sprintf("failed to apply migrations: %v", err))
	}

	fmt.Println("Migrations applied successfully")

}
