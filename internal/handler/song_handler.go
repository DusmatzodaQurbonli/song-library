package handler

import (
	"net/http"
	"strconv"

	"github.com/DusmatzodaQurbonli/song-library/internal/entity"
	"github.com/DusmatzodaQurbonli/song-library/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SongHandler struct {
	service *service.SongService
	logger  *logrus.Logger
}

func NewSongHandler(s *service.SongService, log *logrus.Logger) *SongHandler {
	return &SongHandler{
		service: s,
		logger:  log,
	}
}

// @title Song Library API
// @version 1.0
// @description This is a simple song library API.
// @host localhost:8080
// @BasePath /

// @Summary Get paginated songs
// @Description Get songs with pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param group query string false "Filter by group"
// @Param title query string false "Filter by title"
// @Success 200 {array} entity.Song "List of songs"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	filter := make(map[string]string)
	if group := c.Query("group"); group != "" {
		filter["group"] = group
	}
	if title := c.Query("title"); title != "" {
		filter["title"] = title
	}

	songs, err := h.service.GetSongs(c.Request.Context(), filter, page, size)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get songs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// @Summary Get song text with pagination
// @Description Get song text paginated by verses
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {array} string "Paginated song text"
// @Failure 404 {object} map[string]string "Song not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	verses, err := h.service.GetSongText(c.Request.Context(), uint(id), page, size)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get song text")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verses)
}

// @Summary Add new song
// @Description Add a new song to the library
// @Tags songs
// @Accept json
// @Produce json
// @Param song body entity.Song true "Song Data"
// @Success 201 {object} entity.Song "Created song"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var song entity.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdSong, err := h.service.AddSong(c.Request.Context(), &song)
	if err != nil {
		h.logger.WithError(err).Error("Failed to add song")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdSong)
}

// @Summary Update song
// @Description Update song data
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body entity.Song true "Song Data"
// @Success 200 {object} entity.Song "Updated song"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Song not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var song entity.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	song.ID = uint(id)
	if err := h.service.UpdateSong(c.Request.Context(), &song); err != nil {
		h.logger.WithError(err).Error("Failed to update song")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, song)
}

// @Summary Delete song
// @Description Delete a song by ID
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string "Song not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.DeleteSong(c.Request.Context(), uint(id)); err != nil {
		h.logger.WithError(err).Error("Failed to delete song")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
