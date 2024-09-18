package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/table"
	"github.com/yaninyzwitty/scylla-go-app/configuration"
	"github.com/yaninyzwitty/scylla-go-app/controller"
	"github.com/yaninyzwitty/scylla-go-app/database"
	"github.com/yaninyzwitty/scylla-go-app/repository"
	"github.com/yaninyzwitty/scylla-go-app/router"
	"github.com/yaninyzwitty/scylla-go-app/service"
)

type Person struct {
	ID   gocql.UUID `json:"id"`
	Name string     `json:"name"`
	Age  int        `json:"age"`
}

var personMetadata = table.Metadata{
	Name:    "persons",
	Columns: []string{"id", "name", "age"},
	PartKey: []string{"id"},
}

var session gocqlx.Session

var personTable = table.New(personMetadata)

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

// check this if it will work properly
// func updatePerson(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	if idStr == "" {
// 		http.Error(w, "ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	id, err := gocql.ParseUUID(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
// 		return
// 	}

// 	var person Person
// 	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Construct the update query
// 	query := qb.Update(personTable.Name()).
// 		Set("name", "age").
// 		Where(qb.Eq("id")).Query(session)

// 	// Bind values and execute
// 	err = query.BindMap(qb.M{"id": id, "name": person.Name, "age": person.Age}).ExecRelease()
// 	if err != nil {
// 		http.Error(w, "Failed to update the person", http.StatusInternalServerError)
// 		return
// 	}

// 	responseToJson(w, http.StatusOK, person)
// }
