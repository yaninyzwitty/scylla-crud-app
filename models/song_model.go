package models

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
)

type Song struct {
	ID     gocql.UUID `json:"id"`
	Title  string     `json:"title"`
	Album  string     `json:"album"`
	Artist string     `json:"artist"`
	Tags   []string   `json:"tags"`
	Data   []byte     `json:"data"`
}

var songMetadata = table.Metadata{
	Name:    "songs",
	Columns: []string{"id", "title", "album", "artist", "tags", "data"},
	PartKey: []string{"id"},
}

var SongsTable = table.New(songMetadata)
