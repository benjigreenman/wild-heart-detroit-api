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

type YouTubeProvider struct {
	APIKey    string
	PodcastId string
}

// YouTube API response structs
type YouTubeVideos struct {
	Items []YouTubeVideo `json:"items"`
}

type YouTubeVideo struct {
	ID             string                        `json:"id"`
	Snippet        YouTubeSnippet                `json:"snippet"`
	ContentDetails YouTubeContentDetailsPlaylist `json:"contentDetails"`
}

type YouTubeSnippet struct {
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Thumbnails   YouTubeThumbnails `json:"thumbnails"`
	PublishedAt  string            `json:"publishedAt"`
	ChannelTitle string            `json:"channelTitle"`
	ResourceId   YouTubeResourceId `json:"resourceId"`
}

type YouTubeThumbnails struct {
	Default YouTubeThumbnail `json:"default"`
	Medium  YouTubeThumbnail `json:"medium"`
	High    YouTubeThumbnail `json:"high"`
}

type YouTubeThumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type YouTubeResourceId struct {
	VideoId string `json:"videoId"`
}

type YouTubeContentDetailsPlaylist struct {
	Duration string `json:"duration"`
}

func (y *YouTubeProvider) GetShows(ctx context.Context, limit int) ([]MediaItem, error) {
	baseURL := "https://www.googleapis.com/youtube/v3/playlistItems"
	params := url.Values{}
	params.Set("part", "snippet,contentDetails")
	params.Set("playlistId", y.PodcastId)
	params.Set("key", y.APIKey)
	params.Set("maxResults", strconv.Itoa(limit))

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
			return nil, fmt.Errorf("youtube api error: status %d", resp.StatusCode)
		}
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf(bodyString)
	}

	var apiResp YouTubeVideos
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	// Fetch details for each video concurrently
	enriched, err := y.GetVideoDetails(ctx, apiResp.Items)
	if err != nil {
		return nil, err
	}
	return y.mapEpisodes(enriched), nil
}

func (y *YouTubeProvider) mapEpisodes(eps []YouTubeVideo) []MediaItem {
	items := make([]MediaItem, 0, len(eps))
	for _, ep := range eps {
		items = append(items, MediaItem{
			Title:       ep.Snippet.Title,
			Date:        ExtractDate(ep.Snippet.Description),
			Source:      "youtube",
			Link:        "https://www.youtube.com/embed/" + ep.Snippet.ResourceId.VideoId,
			Description: ep.Snippet.Description,
			Raw:         ep,
			VideoId:     ep.Snippet.ResourceId.VideoId,
			Duration:    ep.ContentDetails.Duration,
		})
	}
	return items
}

// Fetches additional details for each video concurrently
func (y *YouTubeProvider) GetVideoDetails(ctx context.Context, videos []YouTubeVideo) ([]YouTubeVideo, error) {
	type result struct {
		idx  int
		item YouTubeVideo
		err  error
	}
	out := make(chan result, len(videos))

	for i, video := range videos {
		go func(idx int, vid YouTubeVideo) {
			// Fetch video details
			baseURL := "https://www.googleapis.com/youtube/v3/videos"
			params := url.Values{}
			params.Set("part", "snippet,contentDetails")
			params.Set("id", vid.Snippet.ResourceId.VideoId)
			params.Set("key", y.APIKey)
			reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
			if err != nil {
				out <- result{idx, vid, err}
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				out <- result{idx, vid, err}
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				out <- result{idx, vid, fmt.Errorf("youtube video api error: status %d", resp.StatusCode)}
				return
			}
			var wrapper YouTubeVideos
			if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
				out <- result{idx, vid, err}
				return
			}
			if len(wrapper.Items) > 0 {
				vid.ContentDetails.Duration = wrapper.Items[0].ContentDetails.Duration
			}
			out <- result{idx, vid, nil}
		}(i, video)
	}

	videoWithDuration := make([]YouTubeVideo, len(videos))
	var firstErr error
	for i := 0; i < len(videos); i++ {
		res := <-out
		if res.err != nil && firstErr == nil {
			firstErr = res.err
		}
		videoWithDuration[res.idx] = res.item
	}
	close(out)
	if firstErr != nil {
		return nil, firstErr
	}
	return videoWithDuration, nil
}
