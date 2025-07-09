package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	delivery "github.com/bgaurav7/gin-microservice-boilerplate/internal/delivery/http"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/db"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
)

func main() {
	// Check if -dsn flag is provided
	printDSN := flag.Bool("dsn", false, "Print the database connection string and exit")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	
	// If -dsn flag is provided, print the DSN and exit
	if *printDSN {
		fmt.Print(cfg.Database.DSN())
		os.Exit(0)
	}

	// Initialize logger
	log, err := logger.NewLogger(&cfg.Logger)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Initialize database connection
	database, err := db.NewDatabase(&cfg.Database, log)
	if err != nil {
		log.Error("Failed to connect to database", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}

	// Create router
	router := delivery.NewRouter(log, database, cfg)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router.Handler(),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown server
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}
