package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/koccyx/avito_assignment/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		panic("migrations path is required ")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", cfg.Storage.User, cfg.Storage.Password, cfg.Storage.Addres, cfg.Storage.Port, cfg.Storage.Database, cfg.Storage.Schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		cfg.Storage.Database,
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println(err)
	} else {
		fmt.Println("Migrations applied successfully!")
	}
}
