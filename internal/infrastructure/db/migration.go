package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Migrator represents a database migrator
type Migrator struct {
	Config *config.DatabaseConfig
	Logger *logger.Logger
	DB     *Database
}

// NewMigrator creates a new database migrator
func NewMigrator(cfg *config.DatabaseConfig, logger *logger.Logger, db *Database) *Migrator {
	return &Migrator{
		Config: cfg,
		Logger: logger,
		DB:     db,
	}
}

// Run runs database migrations
func (m *Migrator) Run() error {
	m.Logger.Info("Running database migrations", nil)

	// Get SQL DB instance
	sqlDB, err := m.DB.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create migrator with absolute path to migrations
	migrationPath := "file://" + getMigrationsPath()
	m.Logger.Info("Using migration path", map[string]interface{}{"path": migrationPath})
	
	migrator, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Run migrations
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.Logger.Info("Database migrations completed", nil)

	return nil
}

// getMigrationsPath returns the absolute path to the migrations directory
func getMigrationsPath() string {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		// Fallback to relative path if we can't get the working directory
		return "migrations"
	}

	// Check if migrations directory exists in the current working directory
	migrationsPath := filepath.Join(cwd, "migrations")
	if _, err := os.Stat(migrationsPath); err == nil {
		return migrationsPath
	}

	// Check if migrations directory exists in the parent directory (for running from cmd/server)
	parentMigrationsPath := filepath.Join(cwd, "..", "..", "migrations")
	if _, err := os.Stat(parentMigrationsPath); err == nil {
		return parentMigrationsPath
	}

	// Fallback to relative path
	return "migrations"
}
