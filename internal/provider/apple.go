package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type AppleProvider struct {
	PodcastID string
}

type ApplePodcastsSearchResponse struct {
	ResultCount int                   `json:"resultCount"`
	Results     []ApplePodcastEpisode `json:"results"`
}

type ApplePodcastEpisode struct {
	TrackID        int64  `json:"trackId"`
	TrackName      string `json:"trackName"`
	Description    string `json:"description"`
	ReleaseDate    string `json:"releaseDate"`
	ArtworkUrl     string `json:"artworkUrl600"`
	TrackViewUrl   string `json:"trackViewUrl"`
	CollectionName string `json:"collectionName"`
	EpisodeUrl     string `json:"episodeUrl"`
}

func (a *AppleProvider) GetShows(ctx context.Context, limit int) ([]MediaItem, error) {
	baseURL := "https://itunes.apple.com/lookup"
	params := url.Values{}
	params.Set("id", a.PodcastID)
	params.Set("media", "podcast")
	params.Set("entity", "podcastEpisode")
	params.Set("limit", strconv.Itoa(limit))

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("apple api error: status %d", resp.StatusCode)
		}
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf(bodyString)
	}

	var apiResp ApplePodcastsSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return a.mapEpisodes(apiResp.Results), nil
}

func (a *AppleProvider) mapEpisodes(appleEps []ApplePodcastEpisode) []MediaItem {
	items := make([]MediaItem, 0, len(appleEps))
	for _, ep := range appleEps {
		dateString := ExtractDate(ep.Description)
		items = append(items, MediaItem{
			Title:       ep.TrackName,
			Date:        dateString,
			Source:      "apple",
			Link:        ep.TrackViewUrl,
			Description: ep.Description,
			Raw:         ep,
			VideoId:     "",
			Duration:    "",
		})
	}
	return items
}
