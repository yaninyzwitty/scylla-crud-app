package repository

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/qb"
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
	query := qb.Insert(models.SongsTable.Name()).Columns(models.SongsTable.Metadata().Columns...).Query(*r.session)
	err := query.BindStruct(song).ExecRelease()
	if err != nil {
		return models.Song{}, err
	}

	return song, nil

}

func (r *songsRepository) UpdateSong(ctx context.Context, id gocql.UUID, song models.Song) (models.Song, error) {
	query := qb.Update(models.SongsTable.Name()).Set("name", "age").Where(qb.Eq("id")).Query(*r.session) //write manually dont copy all cols...
	err := query.BindStruct(song).ExecRelease()
	if err != nil {
		return models.Song{}, err
	}
	return song, nil

}

func (r *songsRepository) DeleteSong(ctx context.Context, id gocql.UUID) error {
	query := qb.Delete(models.SongsTable.Name()).Where(qb.Eq("id")).Query(*r.session)
	err := query.BindMap(qb.M{"id": id}).ExecRelease()
	if err != nil {
		return err
	}

	return nil
}

func (r *songsRepository) GetAllSongs(ctx context.Context) ([]models.Song, error) {
	return []models.Song{}, nil

}

func (r *songsRepository) GetSong(ctx context.Context, id gocql.UUID) (models.Song, error) {
	var song models.Song
	query := qb.Select(models.SongsTable.Name()).Where(qb.Eq("id")).Query(*r.session)
	query.BindMap(qb.M{"id": id})

	err := query.GetRelease(&song)
	if err != nil {
		return models.Song{}, err
	}

	return song, nil
}
