package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DusmatzodaQurbonli/song-library/internal/config"
	"github.com/DusmatzodaQurbonli/song-library/internal/entity"
	"io"
	"net/http"
	"time"
)

type MusicInfoClient interface {
	GetSongInfo(ctx context.Context, group, title string) (*entity.Song, error)
}

type musicInfoClient struct {
	baseURL string
}

func NewMusicInfoClient(cfg *config.Config) MusicInfoClient {
	return &musicInfoClient{baseURL: cfg.MusicInfoAPI}
}

func (c *musicInfoClient) GetSongInfo(ctx context.Context, group, title string) (*entity.Song, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.baseURL, group, title)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	releaseDate, err := time.Parse("02.01.2006", result.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse release date: %w", err)
	}

	return &entity.Song{
		ReleaseDate: releaseDate,
		Text:        result.Text,
		Link:        result.Link,
	}, nil
}
