package repository

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/yaninyzwitty/scylla-go-app/models"
)

type SongsRepository interface {
	CreateSong(ctx context.Context, song models.Song) (models.Song, error)
	UpdateSong(ctx context.Context, id gocql.UUID, song models.Song) (models.Song, error)
	DeleteSong(ctx context.Context, id gocql.UUID) error
	GetAllSongs(ctx context.Context) ([]models.Song, error)
	GetSong(ctx context.Context, id gocql.UUID) (models.Song, error)
}

type songsRepository struct {
	session *gocqlx.Session
}

func NewSongsRepository(session *gocqlx.Session) SongsRepository {
	return &songsRepository{session: session}
}

func (r *songsRepository) CreateSong(ctx context.Context, song models.Song) (models.Song, error) {
	return models.Song{}, nil

}

func (r *songsRepository) UpdateSong(ctx context.Context, id gocql.UUID, song models.Song) (models.Song, error) {
	return models.Song{}, nil
}

func (r *songsRepository) DeleteSong(ctx context.Context, id gocql.UUID) error {
	return nil
}

func (r *songsRepository) GetAllSongs(ctx context.Context) ([]models.Song, error) {
	return []models.Song{}, nil

}

func (r *songsRepository) GetSong(ctx context.Context, id gocql.UUID) (models.Song, error) {
	return models.Song{}, nil
}
