package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/yaninyzwitty/scylla-go-app/configuration"
	"github.com/yaninyzwitty/scylla-go-app/controller"
	"github.com/yaninyzwitty/scylla-go-app/database"
	"github.com/yaninyzwitty/scylla-go-app/repository"
	"github.com/yaninyzwitty/scylla-go-app/router"
	"github.com/yaninyzwitty/scylla-go-app/service"
)

func main() {
	cfg, err := configuration.NewConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
	}

	session, err := database.NewDatabaseConnection(cfg.HOSTS)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}

	defer session.Close()

	songsRepo := repository.NewSongsRepository(session)
	songsService := service.NewSongsService(songsRepo)
	songsController := controller.NewController(songsService)

	r := router.NewRouter(songsController)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server: %v", err)
		}
	}()

	// Server is running
	slog.Info(fmt.Sprintf("Server is running on port: %s", cfg.PORT))

	// Set up OS signal handling for graceful shutdown
	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, os.Interrupt)

	// Block until signal is received
	<-quitCH
	slog.Info("Received termination signal, shutting down server...")

	// Create context for shutdown timeouts
	shutdownCTX, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCTX); err != nil {
		slog.Error("Failed to gracefully shut down server: %v", err)
	}
	slog.Info("Server shutdown successful")

}
