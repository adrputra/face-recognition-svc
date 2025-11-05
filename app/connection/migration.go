package connection

import (
	"database/sql"
	"face-recognition-svc/app/config"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func MigrateDatabase(c *config.Database) {
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=+07:00",
		c.Host,
		c.Username,
		c.Password,
		c.Database,
		c.Port)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Get the working directory and construct absolute path to migrations
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	migrationsPath := filepath.Join(wd, "migrations")
	migrationsURL := fmt.Sprintf("file://%s", migrationsPath)

	log.Printf("Looking for migrations at: %s", migrationsURL)

	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migration applied successfully!")
}
