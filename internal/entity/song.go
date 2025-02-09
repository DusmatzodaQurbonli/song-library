package entity

import (
	"strings"
	"time"
)

// Song представляет песню в библиотеке
// @Description Song entity
type Song struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Group       string    `gorm:"not null" json:"group"`
	Title       string    `gorm:"not null;index" json:"title"`
	ReleaseDate time.Time `gorm:"not null" json:"release_date"`
	Text        string    `gorm:"type:text;not null" json:"text"`
	Link        string    `gorm:"not null" json:"link"`
}

func (s Song) GetVerses(page, pageSize int) []string {
	verses := strings.Split(s.Text, "\n\n")
	start := (page - 1) * pageSize
	end := page * pageSize

	if start > len(verses) {
		return []string{}
	}
	if end > len(verses) {
		end = len(verses)
	}
	return verses[start:end]
}
