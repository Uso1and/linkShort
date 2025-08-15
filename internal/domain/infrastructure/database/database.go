package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"lnkshrt/internal/domain/config"
)

var DB *sql.DB

func Init() error {
	dbconfig, err := config.LoadConfig()

	if err != nil {
		return fmt.Errorf("failed load .env file")
	}

	configStr := dbconfig.GetConnectionString()

	DB, err = sql.Open("postgres", configStr)

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("success connect database")

	if err := runMigrations(DB); err != nil {
		return fmt.Errorf("failed to run migration:%w", err)
	}

	return nil
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	absPath, err := filepath.Abs("internal/domain/infrastructure/migrations")
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	migrationPath := filepath.ToSlash(absPath)

	if runtime.GOOS == "windows" {

		if len(migrationPath) > 1 && migrationPath[1] == ':' {
			migrationPath = migrationPath[2:]
		}
	}

	log.Println("Migrations path:", migrationPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", absPath)
	}

	files, err := os.ReadDir(absPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}
	log.Printf("Found %d migration files", len(files))

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	log.Printf("Migrations applied successfully. Version: %d, dirty: %v", version, dirty)
	return nil
}
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
