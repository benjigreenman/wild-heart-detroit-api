package provider

import (
	"context"
	"time"
)

type AppleProvider struct {
	APIKey string
}

func (a *AppleProvider) GetShows(ctx context.Context, limit int) ([]MediaItem, error) {
	// TODO: Implement real API call using a.APIKey
	return a.mapEpisodes([]MediaItem{
		{
			ID:          "3",
			Title:       "Apple Episode 1",
			ImageURL:    "https://example.com/image3.jpg",
			ExternalURL: "https://apple.com/ep1",
			DurationMs:  1800000,
			Source:      "apple",
			ReleaseDate: time.Now().AddDate(0, 0, -3),
		},
	})
}

func (a *AppleProvider) mapEpisodes(eps []MediaItem) ([]MediaItem, error) {
	// TODO: Implement real API call using a.APIKey
	return eps, nil
}
