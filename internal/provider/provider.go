package provider

import (
	"context"
	"time"
)

type MediaItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	ImageURL    string    `json:"imageUrl"`
	ExternalURL string    `json:"externalUrl"`
	DurationMs  int64     `json:"durationMs"`
	Source      string    `json:"source"`
	ReleaseDate time.Time `json:"releaseDate"`
}

type Provider interface {
	GetShows(ctx context.Context, limit int) ([]MediaItem, error)
	mapEpisodes(eps []MediaItem) ([]MediaItem, error)
}
