package repository

import (
	"github.com/DusmatzodaQurbonli/song-library/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SongRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewSongRepository(db *gorm.DB, log *logrus.Logger) *SongRepository {
	return &SongRepository{
		db:     db,
		logger: log,
	}
}

func (r *SongRepository) Create(song *entity.Song) error {
	result := r.db.Create(song)
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error": result.Error,
			"group": song.Group,
			"title": song.Title,
		}).Error("Failed to create song")
	}
	return result.Error
}

func (r *SongRepository) GetPaginated(filter map[string]string, page, size int) ([]entity.Song, error) {
	var songs []entity.Song
	query := r.db.Model(&entity.Song{})

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	offset := (page - 1) * size
	err := query.Limit(size).Offset(offset).Find(&songs).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":  err,
			"filter": filter,
			"page":   page,
			"size":   size,
		}).Error("Failed to get songs")
	}

	return songs, err
}

func (r *SongRepository) GetByID(id uint) (*entity.Song, error) {
	var song entity.Song
	err := r.db.First(&song, id).Error

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to get song by ID")
	}

	return &song, err
}

// Update обновляет данные песни
func (r *SongRepository) Update(song *entity.Song) error {
	result := r.db.Save(song)
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error": result.Error,
			"id":    song.ID,
		}).Error("Failed to update song")
	}
	return result.Error
}

// Delete удаляет песню по её ID
func (r *SongRepository) Delete(id uint) error {
	result := r.db.Delete(&entity.Song{}, id)
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error": result.Error,
			"id":    id,
		}).Error("Failed to delete song")
	}
	return result.Error
}
