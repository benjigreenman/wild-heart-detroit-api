package services

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
	"wild-heart-detroit-api/internal/provider"
)

type MediaGroup struct {
	Date          string                          `json:"date"`
	MediaBySource map[string][]provider.MediaItem `json:"mediaBySource"`
}

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

func (m *MediaService) GetAllMediaContent(ctx context.Context, limit int) ([]MediaGroup, error) {
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

func (m *MediaService) groupMediaItemsByDate(items []provider.MediaItem) []MediaGroup {
	// Group items by date
	grouped := make(map[string][]provider.MediaItem)
	for _, item := range items {
		grouped[item.Date] = append(grouped[item.Date], item)
	}

	// Convert map to slice
	groups := make([]MediaGroup, 0, len(grouped))
	for date, dateItems := range grouped {
		groups = append(groups, MediaGroup{Date: date, MediaBySource: m.mapMediaItemsBySource(dateItems)})
	}
	// Sort groups by date descending
	sort.Slice(groups, func(i, j int) bool {
		ti, erri := time.Parse("01/02/2006", groups[i].Date)
		tj, errj := time.Parse("01/02/2006", groups[j].Date)
		if erri != nil || errj != nil {
			return groups[i].Date > groups[j].Date
		}
		return ti.After(tj)
	})

	return groups
}

func (m *MediaService) mapMediaItemsBySource(items []provider.MediaItem) map[string][]provider.MediaItem {
	mapped := make(map[string][]provider.MediaItem)
	for _, item := range items {
		mapped[item.Source] = append(mapped[item.Source], item)
	}
	return mapped
}
