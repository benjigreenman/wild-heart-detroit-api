package provider

import (
	"context"
)

type MediaItem struct {
	Title       string
	Date        string
	Source      string
	Link        string
	Description string
	Raw         any
	VideoId     string
	Duration    string
}

type Provider interface {
	GetShows(ctx context.Context, limit int) ([]MediaItem, error)
}
