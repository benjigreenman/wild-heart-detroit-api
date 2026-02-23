package services

import (
	"context"
	"fmt"
	"sync"
	"wild-heart-detroit-api/internal/provider"
)

type MediaService struct {
	youtube provider.Provider
	spotify provider.Provider
	apple   provider.Provider
}

func NewMediaService(
	y provider.Provider,
	s provider.Provider,
	a provider.Provider,
) *MediaService {
	return &MediaService{
		youtube: y,
		spotify: s,
		apple:   a,
	}
}

func (m *MediaService) GetAllMediaContent(ctx context.Context, limit int) (map[string][]provider.MediaItem, error) {
	Log("MediaService", "Starting GetAllMediaContent with limit="+fmt.Sprint(limit), nil)
	var allItems []provider.MediaItem
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	wg.Add(3)

	go func() {
		defer wg.Done()
		Log("MediaService", "Fetching YouTube shows...", nil)
		items, err := m.youtube.GetShows(ctx, limit)
		if err != nil {
			Error("MediaService", "Error fetching YouTube shows", err)
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
			return
		}
		Log("MediaService", "Got "+fmt.Sprint(len(items))+" YouTube items", nil)
		mu.Lock()
		allItems = append(allItems, items...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		Log("MediaService", "Fetching Spotify shows...", nil)
		items, err := m.spotify.GetShows(ctx, limit)
		if err != nil {
			Error("MediaService", "Error fetching Spotify shows", err)
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
			return
		}
		Log("MediaService", "Got "+fmt.Sprint(len(items))+" Spotify items", nil)
		mu.Lock()
		allItems = append(allItems, items...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		Log("MediaService", "Fetching Apple shows...", nil)
		items, err := m.apple.GetShows(ctx, limit)
		if err != nil {
			Error("MediaService", "Error fetching Apple shows", err)
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
			return
		}
		Log("MediaService", "Got "+fmt.Sprint(len(items))+" Apple items", nil)
		mu.Lock()
		allItems = append(allItems, items...)
		mu.Unlock()
	}()

	wg.Wait()

	if len(errs) > 0 {
		Error("MediaService", "Returning error", errs[0])
		return nil, errs[0] // TODO combine errors
	}

	Log("MediaService", "Returning "+fmt.Sprint(len(allItems))+" total items", nil)
	return m.groupMediaItemsByDate(allItems), nil
}

func (m *MediaService) groupMediaItemsByDate(items []provider.MediaItem) map[string][]provider.MediaItem {
	grouped := make(map[string][]provider.MediaItem)
	for _, item := range items {
		date := item.Date
		grouped[date] = append(grouped[date], item)
	}
	return grouped
}
