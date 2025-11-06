package main

import (
	"backend-challenge/api"
	"backend-challenge/db"
	"backend-challenge/service"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	dbPath := flag.String("db", "data/store.db", "Path to SQLite database")
	flag.Parse()

	// Setup signal-based context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, *port, *dbPath); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, port, dbPath string) error {
	database, router, err := setup(dbPath)
	if err != nil {
		return fmt.Errorf("failed to setup application: %w", err)
	}
	defer database.Close()

	addr := ":" + port
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		fmt.Printf("Server starting on port %s...\n", port)
		fmt.Printf("API available at http://localhost:%s/api\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
	}

	fmt.Println("\nShutting down server gracefully...")

	// Create a deadline for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	fmt.Println("Server stopped")
	return nil
}

// setupApplication initializes database, service, and router
func setup(dbPath string) (*db.DB, http.Handler, error) {
	database, err := db.New(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	svc := service.New(database)
	handler := api.NewHandler(svc)
	router := handler.SetupRoutes()

	return database, router, nil
}
