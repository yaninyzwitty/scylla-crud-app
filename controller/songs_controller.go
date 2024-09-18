package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/yaninyzwitty/scylla-go-app/helpers"
	"github.com/yaninyzwitty/scylla-go-app/models"
	"github.com/yaninyzwitty/scylla-go-app/service"
)

type Controller struct {
	service service.SongsService
}

func NewController(service service.SongsService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) CreateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// generate the song's id
	song.ID = gocql.TimeUUID()

	createdSong, err := c.service.CreateSong(ctx, song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = helpers.NewResponseToJson(w, http.StatusCreated, createdSong)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func (c *Controller) UpdateSong(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var song models.Song
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Id is required!", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// convert the songs id to a uuid (supported)
	id, err := gocql.ParseUUID(idStr)
	if err != nil {
		http.Error(w, "Failed to parse id into uuid!", http.StatusBadRequest)
		return
	}

	updatedSong, err := c.service.UpdateSong(ctx, id, song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = helpers.NewResponseToJson(w, http.StatusOK, updatedSong)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func (c *Controller) DeleteSong(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Id is required!", http.StatusBadRequest)
		return
	}
	// convert the songs id to a uuid (supported)
	id, err := gocql.ParseUUID(idStr)
	if err != nil {
		http.Error(w, "Failed to parse id into uuid!", http.StatusBadRequest)
		return
	}

	err = c.service.DeleteSong(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
func (c *Controller) GetAllSongs(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	// get all songs
	songs, err := c.service.GetAllSongs(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// convert songs to json
	err = helpers.NewResponseToJson(w, http.StatusOK, songs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func (c *Controller) GetSong(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Id is required!", http.StatusBadRequest)
		return
	}
	id, err := gocql.ParseUUID(idStr)
	if err != nil {
		http.Error(w, "Failed to parse id into uuid!", http.StatusBadRequest)
		return
	}

	song, err := c.service.GetSong(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = helpers.NewResponseToJson(w, http.StatusOK, song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

}
