package db

import (
	"fmt"
	"time"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/domain/model"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Database represents a database connection
type Database struct {
	DB     *gorm.DB
	Config *config.DatabaseConfig
	Logger *logger.Logger
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.DatabaseConfig, logger *logger.Logger) (*Database, error) {
	logger.Info("Connecting to database", map[string]interface{}{
		"host": cfg.Host,
		"port": cfg.Port,
		"name": cfg.Name,
	})

	// Create DSN
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name, cfg.SSLMode)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// Create database instance
	database := &Database{
		DB:     db,
		Config: cfg,
		Logger: logger,
	}

	// Auto-migrate models
	if err := database.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	logger.Info("Connected to database", map[string]interface{}{
		"host": cfg.Host,
		"port": cfg.Port,
		"name": cfg.Name,
	})

	return database, nil
}

// AutoMigrate automatically migrates the database schema
func (d *Database) AutoMigrate() error {
	return d.DB.AutoMigrate(&model.Todo{})
}

// Ping checks if the database connection is alive
func (d *Database) Ping() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	return sqlDB.Ping()
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	return sqlDB.Close()
}
