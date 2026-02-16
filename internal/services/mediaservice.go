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

func (m *MediaService) GetAllMediaContent(ctx context.Context, limit int) ([]provider.MediaItem, error) {
	var allItems []provider.MediaItem
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	wg.Add(3)

	go func() {
		defer wg.Done()
		items, err := m.youtube.GetShows(ctx, limit)
		if err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
			return
		}
		mu.Lock()
		allItems = append(allItems, items...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		items, err := m.spotify.GetShows(ctx, limit)
		if err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
			return
		}
		mu.Lock()
		allItems = append(allItems, items...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		items, err := m.apple.GetShows(ctx, limit)
		if err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
			return
		}
		mu.Lock()
		allItems = append(allItems, items...)
		mu.Unlock()
	}()

	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0] // or combine errors
	}

	return allItems, nil
}

func GroupMediaByWeekFromCombined(items []provider.MediaItem) map[string][]provider.MediaItem {
	result := make(map[string][]provider.MediaItem)
	for _, item := range items {
		year, week := item.ReleaseDate.ISOWeek()
		key := fmt.Sprintf("%d-%02d", year, week)
		result[key] = append(result[key], item)
	}
	return result
}
