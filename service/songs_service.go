package service

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/yaninyzwitty/scylla-go-app/models"
	"github.com/yaninyzwitty/scylla-go-app/repository"
)

type SongsService interface {
	CreateSong(ctx context.Context, song models.Song) (models.Song, error)
	UpdateSong(ctx context.Context, id gocql.UUID, song models.Song) (models.Song, error)
	DeleteSong(ctx context.Context, id gocql.UUID) error
	GetAllSongs(ctx context.Context) ([]models.Song, error)
	GetSong(ctx context.Context, id gocql.UUID) (models.Song, error)
}

type songService struct {
	repo repository.SongsRepository
}

func NewSongsService(repo repository.SongsRepository) SongsService {
	return &songService{repo: repo}
}
func (s *songService) CreateSong(ctx context.Context, song models.Song) (models.Song, error) {
	return s.repo.CreateSong(ctx, song)
}

func (s *songService) UpdateSong(ctx context.Context, id gocql.UUID, song models.Song) (models.Song, error) {
	return s.repo.UpdateSong(ctx, id, song)
}

func (s *songService) DeleteSong(ctx context.Context, id gocql.UUID) error {
	return s.repo.DeleteSong(ctx, id)
}

func (s *songService) GetAllSongs(ctx context.Context) ([]models.Song, error) {
	return s.repo.GetAllSongs(ctx)
}

func (s *songService) GetSong(ctx context.Context, id gocql.UUID) (models.Song, error) {
	return s.repo.GetSong(ctx, id)
}
