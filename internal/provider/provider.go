package provider

import (
	"context"
)

type MediaItem struct {
	Title       string `json:"title"`
	Date        string `json:"date"` // mm/dd/yyyy
	Source      string `json:"source"`
	Link        string `json:"link,omitempty"`
	Description string `json:"description,omitempty"`
	Raw         any    `json:"raw,omitempty"`
	VideoId     string `json:"videoId,omitempty"`
	Duration    string `json:"duration,omitempty"`
}

type Provider interface {
	GetShows(ctx context.Context, limit int) ([]MediaItem, error)
}
