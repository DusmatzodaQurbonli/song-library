package service

import (
	"context"
	"github.com/DusmatzodaQurbonli/song-library/internal/entity"
	"github.com/DusmatzodaQurbonli/song-library/internal/repository"
	"github.com/sirupsen/logrus"
)

type SongService struct {
	repo       *repository.SongRepository
	infoClient MusicInfoClient
	logger     *logrus.Logger
}

func NewSongService(repo *repository.SongRepository, client MusicInfoClient, log *logrus.Logger) *SongService {
	return &SongService{
		repo:       repo,
		infoClient: client,
		logger:     log,
	}
}

// AddSong добавляет новую песню в библиотеку
func (s *SongService) AddSong(ctx context.Context, req *entity.Song) (*entity.Song, error) {
	info, err := s.infoClient.GetSongInfo(ctx, req.Group, req.Title)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err,
			"group": req.Group,
			"title": req.Title,
		}).Error("Failed to get song info")
		return nil, err
	}

	req.ReleaseDate = info.ReleaseDate
	req.Text = info.Text
	req.Link = info.Link

	if err := s.repo.Create(req); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err,
			"group": req.Group,
			"title": req.Title,
		}).Error("Failed to create song")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"id":    req.ID,
		"group": req.Group,
		"title": req.Title,
	}).Info("Song created successfully")

	return req, nil
}

// GetSongs возвращает список песен с фильтрацией и пагинацией
func (s *SongService) GetSongs(ctx context.Context, filter map[string]string, page, size int) ([]entity.Song, error) {
	songs, err := s.repo.GetPaginated(filter, page, size)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":  err,
			"filter": filter,
			"page":   page,
			"size":   size,
		}).Error("Failed to get songs")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"count": len(songs),
		"page":  page,
		"size":  size,
	}).Info("Songs retrieved successfully")

	return songs, nil
}

// GetSongText возвращает текст песни с пагинацией по куплетам
func (s *SongService) GetSongText(ctx context.Context, id uint, page, size int) ([]string, error) {
	song, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to get song by ID")
		return nil, err
	}

	verses := song.GetVerses(page, size)
	s.logger.WithFields(logrus.Fields{
		"id":     id,
		"page":   page,
		"size":   size,
		"verses": len(verses),
	}).Info("Song verses retrieved successfully")

	return verses, nil
}

// UpdateSong обновляет данные песни
func (s *SongService) UpdateSong(ctx context.Context, song *entity.Song) error {
	if err := s.repo.Update(song); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err,
			"id":    song.ID,
		}).Error("Failed to update song")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"id": song.ID,
	}).Info("Song updated successfully")

	return nil
}

// DeleteSong удаляет песню по ID
func (s *SongService) DeleteSong(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to delete song")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"id": id,
	}).Info("Song deleted successfully")

	return nil
}
