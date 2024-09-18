package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/scylladb/gocqlx/v3/table"
	"github.com/yaninyzwitty/scylla-go-app/configuration"
	"github.com/yaninyzwitty/scylla-go-app/controller"
	"github.com/yaninyzwitty/scylla-go-app/database"
	"github.com/yaninyzwitty/scylla-go-app/repository"
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

	mux := http.NewServeMux()
	mux.HandleFunc("POST /songs", createSong)        // POST /songs
	mux.HandleFunc("GET /songs/{id}", getSong)       // GET /song/{id}
	mux.HandleFunc("PUT /songs/{id}", updateSong)    // PUT /song/{id}
	mux.HandleFunc("DELETE /songs/{id}", deleteSong) // DELETE /song/{id}

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func createSong(w http.ResponseWriter, r *http.Request) {
	var song Person
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	song.ID = gocql.TimeUUID()

	query := qb.Insert(personTable.Name()).Columns(personTable.Metadata().Columns...).Query(session)
	err := query.BindStruct(song).ExecRelease()
	if err != nil {
		http.Error(w, "Failed to insert the song", http.StatusInternalServerError)
		return
	}
	// json.NewEncoder(w).Encode(song)
	responseToJson(w, http.StatusOK, song)

}

func getSong(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := gocql.ParseUUID(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID format!", http.StatusBadRequest)
		return
	}

	query := qb.Select(personTable.Name()).Where(qb.Eq("id")).Query(session)
	query.BindMap(qb.M{"id": id})

	var song Person
	err = query.GetRelease(&song)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	responseToJson(w, http.StatusOK, song)
	// json.NewEncoder(w).Encode(song)
	// use this instead
}

func updateSong(w http.ResponseWriter, r *http.Request) {
	// Extract ID from the URL path
	idStr := r.PathValue("id")
	fmt.Println(idStr)

	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Parse the UUID from the ID string
	id, err := gocql.ParseUUID(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID format!", http.StatusBadRequest)
		return
	}

	// Decode the request body to get the song data
	var song Person
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the ID of the song to the parsed ID
	song.ID = id

	// Build the update query
	query := qb.Update(personTable.Name()).Set("name", "age").Where(qb.Eq("id")).Query(session) //write manually dont copy all cols...

	// Bind the song struct to the query
	err = query.BindStruct(song).ExecRelease()
	if err != nil {
		http.Error(w, "Failed to update the song", http.StatusInternalServerError)
		return
	}

	// Return the updated song as JSON
	responseToJson(w, http.StatusOK, song)
}

func deleteSong(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := gocql.ParseUUID(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID format!", http.StatusBadRequest)
		return
	}

	query := qb.Delete(personTable.Name()).Where(qb.Eq("id")).Query(session)
	err = query.BindMap(qb.M{"id": id}).ExecRelease()
	if err != nil {
		http.Error(w, "Failed to delete the song", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func responseToJson(w http.ResponseWriter, statusCode int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)

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
