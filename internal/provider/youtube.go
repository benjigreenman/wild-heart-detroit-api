package provider

import (
	"context"
	"time"
)

type YouTubeProvider struct {
	APIKey string
}

func (y *YouTubeProvider) GetShows(ctx context.Context, limit int) ([]MediaItem, error) {
	// TODO: Implement real API call using y.APIKey
	return y.mapEpisodes([]MediaItem{
		{
			ID:          "2",
			Title:       "YouTube Episode 1",
			ImageURL:    "https://example.com/image2.jpg",
			ExternalURL: "https://youtube.com/ep1",
			DurationMs:  2400000,
			Source:      "youtube",
			ReleaseDate: time.Now().AddDate(0, 0, -2),
		},
	})
}

func (y *YouTubeProvider) mapEpisodes(eps []MediaItem) ([]MediaItem, error) {
	// TODO: Implement real API call using y.APIKey
	return eps, nil
}
