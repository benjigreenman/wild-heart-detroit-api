package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SpotifyProvider struct {
	ClientID     string
	ClientSecret string
	PodcastID    string
	AccessToken  string
}

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type SpotifySearchResponse struct {
	Items []SpotifyPodcastEpisode `json:"items"`
}

type SpotifyPodcastEpisode struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ReleaseDate string         `json:"release_date"`
	DurationMs  int64          `json:"duration_ms"`
	ExternalURL string         `json:"external_urls.spotify"`
	Images      []SpotifyImage `json:"images"`
}

type SpotifyImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (s *SpotifyProvider) GetShows(ctx context.Context, limit int) ([]MediaItem, error) {
	_, err := s.getAccessToken(ctx, limit)
	if err != nil {
		return nil, err
	}

	baseURL := "https://api.spotify.com/v1/shows/%s/episodes?limit=%d"
	reqURL := fmt.Sprintf(baseURL, s.PodcastID, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("spotify api error: status %d", resp.StatusCode)
		}
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf(bodyString)
	}

	var apiResp SpotifySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return s.mapEpisodes(apiResp.Items), nil
}

func (s *SpotifyProvider) mapEpisodes(spotifyEps []SpotifyPodcastEpisode) []MediaItem {
	// Map Spotify API episodes to MediaItem
	items := make([]MediaItem, 0, len(spotifyEps))
	for _, ep := range spotifyEps {
		dateString := ExtractDate(ep.Description)
		items = append(items, MediaItem{
			Title:       ep.Name,
			Date:        dateString,
			Source:      "spotify",
			Link:        ep.ExternalURL,
			Description: ep.Description,
			Raw:         ep,
			VideoId:     "",
			Duration:    "",
		})
	}
	return items
}

func (s *SpotifyProvider) getAccessToken(ctx context.Context, limit int) (string, error) {
	reqURL := "https://accounts.spotify.com/api/token"
	body := []byte("grant_type=client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("content-Type", "application/x-www-form-urlencoded")

	credentials := s.ClientID + ":" + s.ClientSecret
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("spotify api error: status %d", resp.StatusCode)
		}
		bodyString := string(bodyBytes)
		return "", fmt.Errorf(bodyString)
	}

	var tokenResponse SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	s.AccessToken = tokenResponse.AccessToken
	return tokenResponse.AccessToken, nil
}
