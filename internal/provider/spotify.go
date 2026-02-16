package provider

import (
	"context"
	"time"
)

type SpotifyProvider struct {
	ClientID     string
	ClientSecret string
}

func (s *SpotifyProvider) GetShows(ctx context.Context, limit int) ([]MediaItem, error) {
	// TODO: Implement real API call using s.ClientID and s.ClientSecret
	return s.mapEpisodes([]MediaItem{
		{
			ID:          "1",
			Title:       "Spotify Episode 1",
			ImageURL:    "https://example.com/image1.jpg",
			ExternalURL: "https://spotify.com/ep1",
			DurationMs:  3600000,
			Source:      "spotify",
			ReleaseDate: time.Now().AddDate(0, 0, -1),
		},
	})
}

func (y *SpotifyProvider) mapEpisodes(eps []MediaItem) ([]MediaItem, error) {
	// TODO: Implement real API call using y.APIKey
	return eps, nil
}
